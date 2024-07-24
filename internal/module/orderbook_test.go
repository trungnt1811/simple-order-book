package module_test

import (
	"testing"
	"time"

	"github.com/trungnt1811/simple-order-book/internal/module"
)

// TestSubmitOrder tests the SubmitOrder function.
func TestOrderBook_SubmitOrder(t *testing.T) {
	// Test case 1: Submit a buy order
	t.Run("Submit Buy Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		customerID := 18
		orderID := orderBook.NextOrderID
		orderBook.SubmitOrder(customerID, 100, true, createGTT(1))

		// Check if the order is added to the BuyOrders heap
		if orderBook.BuyOrders.Len() != 1 {
			t.Errorf("Expected 1 buy order in the heap, got %d", orderBook.BuyOrders.Len())
		}

		// Check if the order is added to the Orders map
		if _, exists := orderBook.Orders[orderID]; !exists {
			t.Errorf("Order ID %d should exist in the Orders map", orderID)
		}

		// Check if the order is added to the CustomerOrders map
		if len(orderBook.CustomerOrders[customerID]) != 1 {
			t.Errorf("Expected 1 order for customer ID %d, got %d", customerID, len(orderBook.CustomerOrders[customerID]))
		}
	})

	// Test case 2: Submit a sell order
	t.Run("Submit Sell Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		customerID := 11
		orderID := orderBook.NextOrderID
		orderBook.SubmitOrder(customerID, 90, false, createGTT(1))

		// Check if the order is added to the SellOrders heap
		if orderBook.SellOrders.Len() != 1 {
			t.Errorf("Expected 1 sell order in the heap, got %d", orderBook.SellOrders.Len())
		}

		// Check if the order is added to the Orders map
		if _, exists := orderBook.Orders[orderID]; !exists {
			t.Errorf("Order ID %d should exist in the Orders map", orderID)
		}

		// Check if the order is added to the CustomerOrders map
		if len(orderBook.CustomerOrders[customerID]) != 1 {
			t.Errorf("Expected 1 order for customer ID %d, got %d", customerID, len(orderBook.CustomerOrders[customerID]))
		}
	})

	// Test case 3: Submit a buy order that matches an existing sell order
	t.Run("Match Buy Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		// Prepare sell order
		sellCustomerID := 1995
		sellOrderID := orderBook.NextOrderID
		orderBook.SubmitOrder(sellCustomerID, 90, false, createGTT(1))

		// Prepare buy order should match
		buyCustomerID := 4953
		buyOrderID := orderBook.NextOrderID
		orderBook.SubmitOrder(buyCustomerID, 90, true, createGTT(1))

		// Check if the buy order is matched and not in the heap
		if orderBook.BuyOrders.Len() != 0 {
			t.Errorf("Expected 0 buy orders in the heap, got %d", orderBook.BuyOrders.Len())
		}

		// Check if the sell order is removed from the heap
		if orderBook.SellOrders.Len() != 0 {
			t.Errorf("Expected 0 sell orders in the heap, got %d", orderBook.SellOrders.Len())
		}

		// Check if the buy order is not added to the heap
		if orderBook.BuyOrders.Len() != 0 {
			t.Errorf("Expected 0 buy orders in the heap, got %d", orderBook.BuyOrders.Len())
		}

		// Check if the matched orders are removed from the Orders map
		if _, exists := orderBook.Orders[sellOrderID]; exists {
			t.Errorf("Order ID %d should not exist in the Orders map", sellOrderID)
		}
		if _, exists := orderBook.Orders[buyOrderID]; exists {
			t.Errorf("Order ID %d should not exist in the Orders map", buyOrderID)
		}

		// Check if the sell order is removed from the CustomerOrders map
		if len(orderBook.CustomerOrders[sellCustomerID]) != 0 {
			t.Errorf("Expected 0 order for customer ID %d, got %d", sellCustomerID, len(orderBook.CustomerOrders[sellCustomerID]))
		}

		// Check if the buy order is not added to the CustomerOrders map
		if len(orderBook.CustomerOrders[buyCustomerID]) != 0 {
			t.Errorf("Expected 0 order for customer ID %d, got %d", buyCustomerID, len(orderBook.CustomerOrders[buyCustomerID]))
		}
	})

	// Test case 4: Submit an order with a nil GTT
	t.Run("Submit Order with Nil GTT", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		orderID := orderBook.NextOrderID
		customerID := 911
		orderBook.SubmitOrder(customerID, 110, true, nil)

		// Check if the order is added to the BuyOrders heap
		if orderBook.BuyOrders.Len() != 1 {
			t.Errorf("Expected 1 buy order in the heap, got %d", orderBook.BuyOrders.Len())
		}

		// Check if the order is added to the Orders map
		if _, exists := orderBook.Orders[orderID]; !exists {
			t.Errorf("Order ID %d should exist in the Orders map", orderID)
		}

		// Check if the order is added to the CustomerOrders map
		if len(orderBook.CustomerOrders[customerID]) != 1 {
			t.Errorf("Expected 1 order for customer ID %d, got %d", customerID, len(orderBook.CustomerOrders[customerID]))
		}
	})

	// Test case 5: Submit multiple orders from the same customer
	t.Run("Submit Multiple Orders from Same Customer", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		customerID := 9947
		orderBook.SubmitOrder(customerID, 120, false, createGTT(2))
		orderBook.SubmitOrder(customerID, 130, false, createGTT(3))

		// Check if the orders are added to the SellOrders heap
		if orderBook.SellOrders.Len() != 2 {
			t.Errorf("Expected 2 sell orders in the heap, got %d", orderBook.SellOrders.Len())
		}

		// Check if the orders are added to the CustomerOrders map
		if len(orderBook.CustomerOrders[customerID]) != 2 {
			t.Errorf("Expected 2 orders for customer ID 1, got %d", len(orderBook.CustomerOrders[customerID]))
		}
	})

	// Test case 6: Submit an order with an expired GTT
	t.Run("Submit Order with Expired GTT", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		expiredGTT := time.Now().Add(-1 * time.Hour)
		orderID := orderBook.NextOrderID
		customerID := 18111995
		orderBook.SubmitOrder(customerID, 95, false, &expiredGTT)

		// Check if the order is not added to the SellOrders heap
		if orderBook.SellOrders.Len() != 0 {
			t.Errorf("Expected 0 sell orders in the heap, got %d", orderBook.SellOrders.Len())
		}

		// Check if the order is not added to the Orders map
		if _, exists := orderBook.Orders[orderID]; exists {
			t.Errorf("Order ID %d should not exist in the Orders map", orderID)
		}

		// Check if the order is not added to the CustomerOrders map
		if len(orderBook.CustomerOrders[customerID]) != 0 {
			t.Errorf("Expected 0 orders for customer ID %d, got %d", customerID, len(orderBook.CustomerOrders[customerID]))
		}
	})
}

// Helper function to create a GTT time.
func createGTT(hours int) *time.Time {
	gtt := time.Now().Add(time.Duration(hours) * time.Hour)
	return &gtt
}
