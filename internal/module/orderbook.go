package module

import (
	"container/heap"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/trungnt1811/simple-order-book/internal/constant"
	"github.com/trungnt1811/simple-order-book/internal/model"
)

// OrderBook manages buy and sell orders.
type OrderBook struct {
	BuyOrders      *model.OrderHeap
	SellOrders     *model.OrderHeap
	Orders         map[uint64]*model.Order          // All orders by ID
	CustomerOrders map[uint]map[uint64]*model.Order // Orders by customer ID and order ID
	NextOrderID    uint64
	mtx            sync.RWMutex
	logger         *zap.Logger
}

// NewOrderBook creates a new OrderBook.
func NewOrderBook(logger *zap.Logger) *OrderBook {
	return &OrderBook{
		BuyOrders:      &model.OrderHeap{Type: constant.BuyOrder},
		SellOrders:     &model.OrderHeap{Type: constant.SellOrder},
		Orders:         make(map[uint64]*model.Order),
		CustomerOrders: make(map[uint]map[uint64]*model.Order),
		NextOrderID:    1,
		logger:         logger,
	}
}

func (ob *OrderBook) SubmitOrder(customerID uint, price uint, orderType constant.OrderType, gtt *time.Time) error {
	ob.mtx.Lock()
	defer ob.mtx.Unlock()

	// Validate inputs
	if orderType != constant.BuyOrder && orderType != constant.SellOrder {
		return fmt.Errorf("invalid order type")
	}

	// Validate price
	if price == 0 {
		return fmt.Errorf("invalid price")
	}

	// Create a new order
	order := &model.Order{
		ID:         ob.NextOrderID,
		CustomerID: customerID,
		Price:      price,
		Timestamp:  time.Now(),
		GTT:        gtt,
		OrderType:  orderType,
	}

	ob.NextOrderID++

	// Try to match the order
	ob.matchOrder(order, orderType)
	return nil
}

// CancelOrder cancels an existing order by ID.
func (ob *OrderBook) CancelOrder(orderID uint64) error {
	ob.mtx.Lock()
	defer ob.mtx.Unlock()

	// Check if the order exists in the order book
	order, exists := ob.Orders[orderID]
	if !exists {
		ob.logger.Debug("Order not found", zap.Uint64("orderID", orderID))
		return fmt.Errorf("order not found: %d", orderID)
	}

	// Remove the order
	ob.removeOrder(order)

	ob.logger.Debug("Order cancelled", zap.Uint64("orderID", orderID))
	return nil
}

// QueryOrders returns all active orders for a given customer ID.
func (ob *OrderBook) QueryOrders(customerID uint) []*model.Order {
	ob.mtx.RLock()
	defer ob.mtx.RUnlock()

	activeOrders := []*model.Order{}
	currentTime := time.Now()

	// Filter and collect only the active orders
	if customerOrders, ok := ob.CustomerOrders[customerID]; ok {
		for _, order := range customerOrders {
			if order.GTT == nil || order.GTT.After(currentTime) {
				activeOrders = append(activeOrders, order)
			}
		}
	}

	return activeOrders // Return the list of active orders for the customer
}

// matchOrder attempts to match a new order with existing orders
func (ob *OrderBook) matchOrder(order *model.Order, orderType constant.OrderType) {
	currentTime := time.Now()

	// If the order's GTT (Good Til Time) is set and it is before the current time, the order is expired.
	// Return immediately as expired orders cannot be matched.
	if order.GTT != nil && order.GTT.Before(currentTime) {
		ob.logger.Debug("Order expired before matching", zap.Uint64("orderID", order.ID))
		return
	}

	var targetOrders, oppositeOrders *model.OrderHeap
	if orderType == constant.BuyOrder {
		targetOrders = ob.BuyOrders
		oppositeOrders = ob.SellOrders
	} else {
		targetOrders = ob.SellOrders
		oppositeOrders = ob.BuyOrders
	}

	skippedOrders := []*model.Order{}

	// Attempt to match the order with existing opposite orders
	for oppositeOrders.Len() > 0 {
		// Retrieve the top opposite order
		oppositeOrder := heap.Pop(oppositeOrders).(*model.Order)

		// Skip if the opposite order is not present in the Orders map
		if _, exists := ob.Orders[oppositeOrder.ID]; !exists {
			continue
		}

		// Skip if the opposite order belongs to the same customer
		if oppositeOrder.CustomerID == order.CustomerID {
			skippedOrders = append(skippedOrders, oppositeOrder) // Temporarily store the order
			continue
		}

		// Check if the opposite order is still active based on its GTT (Good Til Time)
		if oppositeOrder.GTT == nil || oppositeOrder.GTT.After(currentTime) {
			// Check if the order prices can match
			if (orderType == constant.BuyOrder && oppositeOrder.Price <= order.Price) ||
				(orderType == constant.SellOrder && oppositeOrder.Price >= order.Price) {
				// A match is found, execute the trade
				ob.logger.Info("Matched orders",
					zap.Uint64("orderID", order.ID),
					zap.Uint64("oppositeOrderID", oppositeOrder.ID),
					zap.Uint("price", oppositeOrder.Price),
				)

				// Remove the matched opposite order
				ob.removeOrder(oppositeOrder)

				// Reinsert any skipped orders before returning
				ob.reinsertSkippedOrders(oppositeOrders, skippedOrders)
				return // Exit after successful match
			}
		} else {
			// Remove expired opposite orders
			ob.removeOrder(oppositeOrder)
			continue
		}

		// No match found, push the opposite order back and exit loop
		heap.Push(oppositeOrders, oppositeOrder)
		// Reinsert any skipped orders before exiting the loop
		ob.reinsertSkippedOrders(oppositeOrders, skippedOrders)
		break
	}

	// No match found, add the order to the list of active target orders
	heap.Push(targetOrders, order)
	ob.Orders[order.ID] = order

	// Add the order to the CustomerOrders map
	if ob.CustomerOrders[order.CustomerID] == nil {
		ob.CustomerOrders[order.CustomerID] = make(map[uint64]*model.Order)
	}
	ob.CustomerOrders[order.CustomerID][order.ID] = order
}

// removeOrder remove an order from all relevant data structures
func (ob *OrderBook) removeOrder(order *model.Order) {
	delete(ob.Orders, order.ID)
	if customerOrders, ok := ob.CustomerOrders[order.CustomerID]; ok {
		delete(customerOrders, order.ID)
		if len(customerOrders) == 0 {
			delete(ob.CustomerOrders, order.CustomerID)
		}
	}
}

// reinsertSkippedOrders reinsert skipped orders back into the heap
func (ob *OrderBook) reinsertSkippedOrders(orders *model.OrderHeap, skippedOrders []*model.Order) {
	for _, skipped := range skippedOrders {
		heap.Push(orders, skipped)
	}
}
