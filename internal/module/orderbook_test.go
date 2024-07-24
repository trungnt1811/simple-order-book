package module_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trungnt1811/simple-order-book/internal/module"
)

// Helper function to create a GTT time.
func createGTT(hours int) *time.Time {
	gtt := time.Now().Add(time.Duration(hours) * time.Hour)
	return &gtt
}

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
		require.Equal(t, 1, orderBook.BuyOrders.Len(), "Expected 1 buy order in the heap")

		// Check if the order is added to the Orders map
		_, exists := orderBook.Orders[orderID]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID)

		// Check if the order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.CustomerOrders[customerID]), "Expected 1 order for customer ID %d", customerID)
	})

	// Test case 2: Submit a sell order
	t.Run("Submit Sell Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		customerID := 11
		orderID := orderBook.NextOrderID
		orderBook.SubmitOrder(customerID, 90, false, createGTT(1))

		// Check if the order is added to the SellOrders heap
		require.Equal(t, 1, orderBook.SellOrders.Len(), "Expected 1 sell order in the heap")

		// Check if the order is added to the Orders map
		_, exists := orderBook.Orders[orderID]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID)

		// Check if the order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.CustomerOrders[customerID]), "Expected 1 order for customer ID %d", customerID)
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
		require.Equal(t, 0, orderBook.BuyOrders.Len(), "Expected 0 buy orders in the heap")

		// Check if the sell order is removed from the heap
		require.Equal(t, 0, orderBook.SellOrders.Len(), "Expected 0 sell orders in the heap")

		// Check if the matched orders are removed from the Orders map
		_, sellOrderExists := orderBook.Orders[sellOrderID]
		require.False(t, sellOrderExists, "Order ID %d should not exist in the Orders map", sellOrderID)
		_, buyOrderExists := orderBook.Orders[buyOrderID]
		require.False(t, buyOrderExists, "Order ID %d should not exist in the Orders map", buyOrderID)

		// Check if the sell order is removed from the CustomerOrders map
		require.Equal(t, 0, len(orderBook.CustomerOrders[sellCustomerID]), "Expected 0 order for customer ID %d", sellCustomerID)

		// Check if the buy order is not added to the CustomerOrders map
		require.Equal(t, 0, len(orderBook.CustomerOrders[buyCustomerID]), "Expected 0 order for customer ID %d", buyCustomerID)
	})

	// Test case 4: Submit an order with a nil GTT
	t.Run("Submit Order with Nil GTT", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		orderID := orderBook.NextOrderID
		customerID := 911
		orderBook.SubmitOrder(customerID, 110, true, nil)

		// Check if the order is added to the BuyOrders heap
		require.Equal(t, 1, orderBook.BuyOrders.Len(), "Expected 1 buy order in the heap")

		// Check if the order is added to the Orders map
		_, exists := orderBook.Orders[orderID]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID)

		// Check if the order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.CustomerOrders[customerID]), "Expected 1 order for customer ID %d", customerID)
	})

	// Test case 5: Submit multiple orders from the same customer
	t.Run("Submit Multiple Orders from Same Customer", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		customerID := 9947
		orderBook.SubmitOrder(customerID, 120, false, createGTT(2))
		orderBook.SubmitOrder(customerID, 130, false, createGTT(3))

		// Check if the orders are added to the SellOrders heap
		require.Equal(t, 2, orderBook.SellOrders.Len(), "Expected 2 sell orders in the heap")

		// Check if the orders are added to the CustomerOrders map
		require.Equal(t, 2, len(orderBook.CustomerOrders[customerID]), "Expected 2 orders for customer ID %d", customerID)
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
		require.Equal(t, 0, orderBook.SellOrders.Len(), "Expected 0 sell orders in the heap")

		// Check if the order is not added to the Orders map
		_, exists := orderBook.Orders[orderID]
		require.False(t, exists, "Order ID %d should not exist in the Orders map", orderID)

		// Check if the order is not added to the CustomerOrders map
		require.Equal(t, 0, len(orderBook.CustomerOrders[customerID]), "Expected 0 orders for customer ID %d", customerID)
	})

	// Test case 9: Submit Orders from Same Customer with Same Prices but Different Timestamp
	t.Run("Submit Orders from Same Customer with Same Prices but Different Timestamp", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		customerID := 69
		orderID1 := orderBook.NextOrderID
		orderBook.SubmitOrder(customerID, 100, true, createGTT(1))

		// Change the timestamp but same price
		orderID2 := orderBook.NextOrderID
		orderBook.SubmitOrder(customerID, 100, true, createGTT(2))

		// Check if both orders are added to the BuyOrders heap
		require.Equal(t, 2, orderBook.BuyOrders.Len(), "Expected 2 buy orders in the heap")

		// Check if both orders are added to the Orders map
		_, exists1 := orderBook.Orders[orderID1]
		require.True(t, exists1, "Order ID %d should exist in the Orders map", orderID1)
		_, exists2 := orderBook.Orders[orderID2]
		require.True(t, exists2, "Order ID %d should exist in the Orders map", orderID2)

		// Check if both orders are added to the CustomerOrders map
		require.Equal(t, 2, len(orderBook.CustomerOrders[customerID]), "Expected 2 orders for customer ID %d", customerID)
	})
}
