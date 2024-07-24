package module

import (
	"container/heap"
	"fmt"
	"sync"
	"time"

	"github.com/trungnt1811/simple-order-book/internal/model"
)

// OrderBook manages buy and sell orders.
type OrderBook struct {
	BuyOrders      *model.OrderHeap
	SellOrders     *model.OrderHeap
	Orders         map[int]*model.Order   // All orders by ID
	CustomerOrders map[int][]*model.Order // Orders by customer ID
	NextOrderID    int
	mu             sync.Mutex
}

// NewOrderBook creates a new OrderBook.
func NewOrderBook() *OrderBook {
	return &OrderBook{
		BuyOrders:      &model.OrderHeap{Desc: true},
		SellOrders:     &model.OrderHeap{Desc: false},
		Orders:         make(map[int]*model.Order),
		CustomerOrders: make(map[int][]*model.Order),
		NextOrderID:    1,
	}
}

// SubmitOrder submits a new buy or sell order.
func (ob *OrderBook) SubmitOrder(customerID int, price int, isBuy bool, gtt *time.Time) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	// Create a new order
	order := &model.Order{
		ID:         ob.NextOrderID,
		CustomerID: customerID,
		Price:      price,
		Timestamp:  time.Now(),
		GTT:        gtt,
	}

	ob.NextOrderID++

	// Try to match the order
	if isBuy {
		ob.matchBuyOrder(order)
	} else {
		ob.matchSellOrder(order)
	}
}

// CancelOrder cancels an existing order by ID.
func (ob *OrderBook) CancelOrder(orderID int) {
	ob.mu.Lock()         // Lock the order book for safe concurrent access
	defer ob.mu.Unlock() // Ensure the lock is released when the function returns

	// Check if the order exists in the order book
	order, exists := ob.Orders[orderID]
	if !exists {
		fmt.Println("Order not found") // Inform the user if the order is not found
		return
	}

	// Remove the order from the main order book
	delete(ob.Orders, orderID)

	// Remove the order from the customer's list of orders
	customerOrders := ob.CustomerOrders[order.CustomerID]
	for i, o := range customerOrders {
		if o.ID == orderID {
			// Remove the order from the slice of the customer's orders
			ob.CustomerOrders[order.CustomerID] = append(customerOrders[:i], customerOrders[i+1:]...)
			break
		}
	}

	fmt.Printf("Order %d cancelled\n", orderID) // Confirm the cancellation
}

// QueryOrders returns all active orders for a given customer ID.
func (ob *OrderBook) QueryOrders(customerID int) []*model.Order {
	ob.mu.Lock()         // Lock the order book to prevent concurrent modifications
	defer ob.mu.Unlock() // Ensure the lock is released after the function completes

	// Retrieve the orders for the given customer ID
	customerOrders := ob.CustomerOrders[customerID]
	activeOrders := []*model.Order{}
	currentTime := time.Now()

	// Filter and collect only the active orders
	for _, order := range customerOrders {
		// Check if the order is still active based on its GTT (Good Til Time)
		if order.GTT == nil || order.GTT.After(currentTime) {
			activeOrders = append(activeOrders, order)
		}
	}

	return activeOrders // Return the list of active orders for the customer
}

func (ob *OrderBook) matchBuyOrder(order *model.Order) {
	currentTime := time.Now()

	// If the buy order's GTT (Good Til Time) is set and it is before the current time, the order is expired.
	// Return immediately as expired orders cannot be matched.
	if order.GTT != nil && order.GTT.Before(currentTime) {
		return
	}

	skippedOrders := []*model.Order{}

	// Attempt to match the buy order with existing sell orders
	for ob.SellOrders.Len() > 0 {
		// Retrieve the lowest sell order
		sellOrder := heap.Pop(ob.SellOrders).(*model.Order)

		// Skip if the sell order belongs to the same customer
		if sellOrder.CustomerID == order.CustomerID {
			skippedOrders = append(skippedOrders, sellOrder) // Temporarily store the order
			continue
		}

		// Check if the order is still active based on its GTT (Good Til Time)
		if sellOrder.GTT == nil || sellOrder.GTT.After(currentTime) {
			// Check if the buy price can match the sell price
			if sellOrder.Price <= order.Price {
				// A match is found, execute the trade
				fmt.Printf("Matched Buy Order %d with Sell Order %d at price %d\n", order.ID, sellOrder.ID, sellOrder.Price)

				// Remove the matched sell order
				delete(ob.Orders, sellOrder.ID)

				// Remove the matched sell order from the customer's orders
				ob.removeOrderFromCustomerOrders(sellOrder.CustomerID, sellOrder.ID)

				// Reinsert any skipped orders before returning
				for _, skipped := range skippedOrders {
					heap.Push(ob.SellOrders, skipped)
				}
				return // Exit after successful match
			}
		} else {
			delete(ob.Orders, sellOrder.ID)                                      // Remove expired sell orders
			ob.removeOrderFromCustomerOrders(sellOrder.CustomerID, sellOrder.ID) // Remove the expired sell order from the customer's orders
			continue
		}

		// No match found, push the sell order back and exit loop
		heap.Push(ob.SellOrders, sellOrder)
		// Reinsert any skipped orders before exiting the loop
		for _, skipped := range skippedOrders {
			heap.Push(ob.SellOrders, skipped)
		}
		break
	}

	// No match found, add the buy order to the list of active buy orders
	heap.Push(ob.BuyOrders, order)
	ob.Orders[order.ID] = order
	ob.CustomerOrders[order.CustomerID] = append(ob.CustomerOrders[order.CustomerID], order)
}

func (ob *OrderBook) matchSellOrder(order *model.Order) {
	currentTime := time.Now()

	// If the sell order's GTT (Good Til Time) is set and it is before the current time, the order is expired.
	// Return immediately as expired orders cannot be matched.
	if order.GTT != nil && order.GTT.Before(currentTime) {
		return
	}

	skippedOrders := []*model.Order{}

	// Attempt to match the sell order with existing buy orders
	for ob.BuyOrders.Len() > 0 {
		// Retrieve the highest buy order
		buyOrder := heap.Pop(ob.BuyOrders).(*model.Order)

		// Skip if the buy order belongs to the same customer
		if buyOrder.CustomerID == order.CustomerID {
			skippedOrders = append(skippedOrders, buyOrder) // Temporarily store the order
			continue
		}

		// Check if the order is still active based on its GTT (Good Til Time)
		if buyOrder.GTT == nil || buyOrder.GTT.After(currentTime) {
			// Check if the sell price can match the buy price
			if buyOrder.Price >= order.Price {
				// A match is found, execute the trade
				fmt.Printf("Matched Sell Order %d with Buy Order %d at price %d\n", order.ID, buyOrder.ID, buyOrder.Price)

				// Remove the matched buy order
				delete(ob.Orders, buyOrder.ID)

				// Remove the matched buy order from the customer's orders
				ob.removeOrderFromCustomerOrders(buyOrder.CustomerID, buyOrder.ID)

				// Reinsert any skipped orders before returning
				for _, skipped := range skippedOrders {
					heap.Push(ob.BuyOrders, skipped)
				}
				return // Exit after successful match
			}
		} else {
			delete(ob.Orders, buyOrder.ID)                                     // Remove expired buy orders
			ob.removeOrderFromCustomerOrders(buyOrder.CustomerID, buyOrder.ID) // Remove the expired buy order from the customer's orders
			continue
		}

		// No match found, push the buy order back and exit loop
		heap.Push(ob.BuyOrders, buyOrder)
		// Reinsert any skipped orders before exiting the loop
		for _, skipped := range skippedOrders {
			heap.Push(ob.BuyOrders, skipped)
		}
		break
	}

	// No match found, add the sell order to the list of active sell orders
	heap.Push(ob.SellOrders, order)
	ob.Orders[order.ID] = order
	ob.CustomerOrders[order.CustomerID] = append(ob.CustomerOrders[order.CustomerID], order)
}

// Helper function to remove an order from the CustomerOrders map
func (ob *OrderBook) removeOrderFromCustomerOrders(customerID, orderID int) {
	customerOrders := ob.CustomerOrders[customerID]
	for i, o := range customerOrders {
		if o.ID == orderID {
			// Remove the order from the slice of the customer's orders
			ob.CustomerOrders[customerID] = append(customerOrders[:i], customerOrders[i+1:]...)
			break
		}
	}
}
