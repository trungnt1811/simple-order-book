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

	// Attempt to match the buy order with existing sell orders
	for ob.SellOrders.Len() > 0 {
		// Retrieve the lowest sell order
		sellOrder := heap.Pop(ob.SellOrders).(*model.Order)

		// Skip if the sell order belongs to the same customer
		if sellOrder.CustomerID == order.CustomerID {
			heap.Push(ob.SellOrders, sellOrder) // Push it back to the heap
			continue
		}

		// Check if the order is still active based on its GTT (Good Til Time)
		if sellOrder.GTT == nil || sellOrder.GTT.After(currentTime) {
			// Check if the buy price can match the sell price
			if sellOrder.Price <= order.Price {
				// A match is found, execute the trade
				fmt.Printf("Matched Buy Order %d with Sell Order %d at price %d\n", order.ID, sellOrder.ID, sellOrder.Price)
				delete(ob.Orders, sellOrder.ID) // Remove the matched sell order
				return                          // Exit after successful match
			}
		} else {
			delete(ob.Orders, sellOrder.ID) // Remove expired sell orders
			continue
		}

		// No match found, push the sell order back and exit loop
		heap.Push(ob.SellOrders, sellOrder)
		break
	}

	// No match found, add the buy order to the list of active buy orders
	heap.Push(ob.BuyOrders, order)
	ob.Orders[order.ID] = order
	ob.CustomerOrders[order.CustomerID] = append(ob.CustomerOrders[order.CustomerID], order)
}

func (ob *OrderBook) matchSellOrder(order *model.Order) {
	currentTime := time.Now()

	// Attempt to match the sell order with existing buy orders
	for ob.BuyOrders.Len() > 0 {
		// Retrieve the highest buy order
		buyOrder := heap.Pop(ob.BuyOrders).(*model.Order)

		// Skip if the buy order belongs to the same customer
		if buyOrder.CustomerID == order.CustomerID {
			heap.Push(ob.BuyOrders, buyOrder) // Push it back to the heap
			continue
		}

		// Check if the order is still active based on its GTT (Good Til Time)
		if buyOrder.GTT == nil || buyOrder.GTT.After(currentTime) {
			// Check if the sell price can match the buy price
			if buyOrder.Price >= order.Price {
				// A match is found, execute the trade
				fmt.Printf("Matched Sell Order %d with Buy Order %d at price %d\n", order.ID, buyOrder.ID, buyOrder.Price)
				delete(ob.Orders, buyOrder.ID) // Remove the matched buy order
				return                         // Exit after successful match
			}
		} else {
			delete(ob.Orders, buyOrder.ID) // Remove expired buy orders
			continue
		}

		// No match found, push the buy order back and exit loop
		heap.Push(ob.BuyOrders, buyOrder)
		break
	}

	// No match found, add the sell order to the list of active sell orders
	heap.Push(ob.SellOrders, order)
	ob.Orders[order.ID] = order
	ob.CustomerOrders[order.CustomerID] = append(ob.CustomerOrders[order.CustomerID], order)
}
