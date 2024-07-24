package module_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trungnt1811/simple-order-book/internal/constant"
	"github.com/trungnt1811/simple-order-book/internal/module"
)

// Helper function to create a GTT time.
func createGTT(hours int) *time.Time {
	gtt := time.Now().Add(time.Duration(hours) * time.Hour)
	return &gtt
}

// TestSubmitOrder tests the SubmitOrder function.
func TestOrderBook_SubmitOrder(t *testing.T) {
	t.Run("Submit Buy Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		customerID := uint(18)
		orderID := orderBook.NextOrderID
		orderBook.SubmitOrder(customerID, 100, constant.BuyOrder, createGTT(1))

		// Check if the order is added to the BuyOrders heap
		require.Equal(t, 1, orderBook.BuyOrders.Len(), "Expected 1 buy order in the heap")

		// Check if the order is added to the Orders map
		_, exists := orderBook.Orders[orderID]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID)

		// Check if the order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.CustomerOrders[customerID]), "Expected 1 order for customer ID %d", customerID)
	})

	t.Run("Submit Sell Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		customerID := uint(11)
		orderID := orderBook.NextOrderID
		orderBook.SubmitOrder(customerID, 90, constant.SellOrder, createGTT(1))

		// Check if the order is added to the SellOrders heap
		require.Equal(t, 1, orderBook.SellOrders.Len(), "Expected 1 sell order in the heap")

		// Check if the order is added to the Orders map
		_, exists := orderBook.Orders[orderID]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID)

		// Check if the order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.CustomerOrders[customerID]), "Expected 1 order for customer ID %d", customerID)
	})

	t.Run("Match Buy Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		// Prepare sell order
		sellCustomerID := uint(1995)
		sellOrderID := orderBook.NextOrderID
		orderBook.SubmitOrder(sellCustomerID, 90, constant.SellOrder, createGTT(1))

		// Prepare buy order should match
		buyCustomerID := uint(4953)
		buyOrderID := orderBook.NextOrderID
		orderBook.SubmitOrder(buyCustomerID, 90, constant.BuyOrder, createGTT(1))

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

	t.Run("Submit Order with Nil GTT", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		orderID := orderBook.NextOrderID
		customerID := uint(911)
		orderBook.SubmitOrder(customerID, 110, constant.BuyOrder, nil)

		// Check if the order is added to the BuyOrders heap
		require.Equal(t, 1, orderBook.BuyOrders.Len(), "Expected 1 buy order in the heap")

		// Check if the order is added to the Orders map
		_, exists := orderBook.Orders[orderID]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID)

		// Check if the order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.CustomerOrders[customerID]), "Expected 1 order for customer ID %d", customerID)
	})

	t.Run("Submit Multiple Orders from Same Customer", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		customerID := uint(9947)
		orderBook.SubmitOrder(customerID, 120, constant.SellOrder, createGTT(2))
		orderBook.SubmitOrder(customerID, 130, constant.SellOrder, createGTT(3))

		// Check if the orders are added to the SellOrders heap
		require.Equal(t, 2, orderBook.SellOrders.Len(), "Expected 2 sell orders in the heap")

		// Check if the orders are added to the CustomerOrders map
		require.Equal(t, 2, len(orderBook.CustomerOrders[customerID]), "Expected 2 orders for customer ID %d", customerID)
	})

	t.Run("Submit Order with Expired GTT", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		expiredGTT := time.Now().Add(-1 * time.Hour)
		orderID := orderBook.NextOrderID
		customerID := uint(18111995)
		orderBook.SubmitOrder(customerID, 95, constant.SellOrder, &expiredGTT)

		// Check if the order is not added to the SellOrders heap
		require.Equal(t, 0, orderBook.SellOrders.Len(), "Expected 0 sell orders in the heap")

		// Check if the order is not added to the Orders map
		_, exists := orderBook.Orders[orderID]
		require.False(t, exists, "Order ID %d should not exist in the Orders map", orderID)

		// Check if the order is not added to the CustomerOrders map
		require.Equal(t, 0, len(orderBook.CustomerOrders[customerID]), "Expected 0 orders for customer ID %d", customerID)
	})

	t.Run("Submit Orders from Same Customer with Same Prices but Different Timestamp", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		customerID := uint(69)
		orderID1 := orderBook.NextOrderID
		orderBook.SubmitOrder(customerID, 100, constant.BuyOrder, createGTT(1))

		// Change the timestamp but same price
		orderID2 := orderBook.NextOrderID
		orderBook.SubmitOrder(customerID, 100, constant.BuyOrder, createGTT(2))

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

	t.Run("Match Cancelled Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		// Submit and then cancel a buy order
		buyCustomerID := uint(3456)
		buyOrderID := orderBook.NextOrderID
		orderBook.SubmitOrder(buyCustomerID, 95, constant.BuyOrder, createGTT(1))
		orderBook.CancelOrder(buyOrderID)

		// Submit a sell order with matching price
		sellCustomerID := uint(7890)
		sellOrderID := orderBook.NextOrderID
		orderBook.SubmitOrder(sellCustomerID, 95, constant.SellOrder, createGTT(1))

		// Check if the cancelled buy order is not matched
		require.Equal(t, 0, orderBook.BuyOrders.Len(), "Expected 0 buy orders in the heap after cancellation")
		require.Equal(t, 1, orderBook.SellOrders.Len(), "Expected 1 sell order in the heap after attempting to match with cancelled buy order")

		// Check if the sell order is still present in the Orders map
		_, sellOrderExists := orderBook.Orders[sellOrderID]
		require.True(t, sellOrderExists, "Sell order ID %d should exist in the Orders map", sellOrderID)

		// Check if the cancelled buy order is not in the Orders map
		_, buyOrderExists := orderBook.Orders[buyOrderID]
		require.False(t, buyOrderExists, "Cancelled buy order ID %d should not exist in the Orders map", buyOrderID)

		// Check if the sell order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.CustomerOrders[sellCustomerID]), "Expected 1 order for customer ID %d", sellCustomerID)

		// Check if the cancelled buy order is not in the CustomerOrders map
		require.Equal(t, 0, len(orderBook.CustomerOrders[buyCustomerID]), "Expected 0 orders for customer ID %d", buyCustomerID)
	})
}

func TestOrderBook_CancelOrder(t *testing.T) {
	t.Run("Cancel existing order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		orderID := orderBook.NextOrderID

		// Submit initial order
		customerID := uint(123)
		orderBook.SubmitOrder(customerID, 100, constant.BuyOrder, createGTT(1))

		// Check if the order is added to the Orders map
		_, exists := orderBook.Orders[orderID]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID)

		// Check if the order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.CustomerOrders[customerID]), "Unexpected number of orders for customer ID %d", customerID)

		// Perform cancel existing order
		err := orderBook.CancelOrder(orderID)
		require.NoError(t, err, "CancelOrder should not return an error")

		// Check if the order is removed from the Orders map
		_, exists = orderBook.Orders[orderID]
		require.False(t, exists, "Order ID %d should not exist in the Orders map", orderID)

		// Check if the order is removed from the CustomerOrders map
		require.Equal(t, 0, len(orderBook.CustomerOrders[customerID]), "Unexpected number of orders for customer ID %d", customerID)
	})
	t.Run("Cancel non-existent order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		orderID := orderBook.NextOrderID

		// Submit initial order
		customerID := uint(456)
		orderBook.SubmitOrder(customerID, 100, constant.BuyOrder, createGTT(1))

		// Check if the order is added to the Orders map
		_, exists := orderBook.Orders[orderID]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID)

		// Check if the order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.CustomerOrders[customerID]), "Unexpected number of orders for customer ID %d", customerID)

		// Try to cancel a non-existent order
		nonExistentOrderID := orderID + 1
		err := orderBook.CancelOrder(nonExistentOrderID)
		require.Error(t, err, fmt.Sprintf("order not found: %d", nonExistentOrderID))

		// Check if the order is still present in the Orders map
		_, exists = orderBook.Orders[orderID]
		require.True(t, exists, "Order ID %d should still exist in the Orders map", orderID)

		// Check if the order is still present in the CustomerOrders map
		require.Equal(t, 1, len(orderBook.CustomerOrders[customerID]), "Unexpected number of orders for customer ID %d", customerID)
	})

	t.Run("Cancel order in multi-order customer", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBook()

		// Submit initial orders
		customerID := uint(789)
		orderID1 := orderBook.NextOrderID
		orderBook.SubmitOrder(customerID, 150, constant.BuyOrder, createGTT(1))
		orderID2 := orderBook.NextOrderID
		orderBook.SubmitOrder(customerID, 160, constant.BuyOrder, createGTT(1))

		// Cancel the first order
		err := orderBook.CancelOrder(orderID1)
		require.NoError(t, err, "CancelOrder should not return an error")

		// Check if the first order is removed from the Orders map
		_, exists := orderBook.Orders[orderID1]
		require.False(t, exists, "Order ID %d should not exist in the Orders map", orderID1)

		// Check if the second order is still present in the Orders map
		_, exists = orderBook.Orders[orderID2]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID2)

		// Check if the CustomerOrders map is updated correctly
		require.Equal(t, 1, len(orderBook.CustomerOrders[customerID]), "Unexpected number of orders for customer ID %d", customerID)
	})
}
