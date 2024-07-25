package module_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/trungnt1811/simple-order-book/internal/constant"
	"github.com/trungnt1811/simple-order-book/internal/module"
	"github.com/trungnt1811/simple-order-book/internal/util"
)

// TestSubmitOrder tests the SubmitOrder function.
func TestOrderBookUCase_SubmitOrder(t *testing.T) {
	logger := util.SetupLogger()
	defer logger.Sync() // Flushes buffer, if any
	t.Run("Submit Buy Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		customerID := uint(18)
		orderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(customerID, 100, constant.BuyOrder, util.CreateGTT(1))

		// Check if the order is added to the BuyOrders heap
		require.Equal(t, 1, orderBook.GetBuyOrders().Len(), "Expected 1 buy order in the heap")

		// Check if the order is added to the Orders map
		_, exists := orderBook.GetOrders()[orderID]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID)

		// Check if the order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.GetCustomerOrders()[customerID]), "Expected 1 order for customer ID %d", customerID)
	})

	t.Run("Submit Sell Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		customerID := uint(11)
		orderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(customerID, 90, constant.SellOrder, util.CreateGTT(1))

		// Check if the order is added to the SellOrders heap
		require.Equal(t, 1, orderBook.GetSellOrders().Len(), "Expected 1 sell order in the heap")

		// Check if the order is added to the Orders map
		_, exists := orderBook.GetOrders()[orderID]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID)

		// Check if the order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.GetCustomerOrders()[customerID]), "Expected 1 order for customer ID %d", customerID)
	})

	t.Run("Submit Buy Order with Exact Price and Match Sell Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		// Prepare sell order
		sellCustomerID := uint(1995)
		sellOrderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(sellCustomerID, 90, constant.SellOrder, util.CreateGTT(1))

		// Prepare buy order should match
		buyCustomerID := uint(4953)
		buyOrderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(buyCustomerID, 90, constant.BuyOrder, util.CreateGTT(1))

		// Check if the buy order is matched and not in the heap
		require.Equal(t, 0, orderBook.GetBuyOrders().Len(), "Expected 0 buy orders in the heap")

		// Check if the sell order is removed from the heap
		require.Equal(t, 0, orderBook.GetSellOrders().Len(), "Expected 0 sell orders in the heap")

		// Check if the matched orders are removed from the Orders map
		_, sellOrderExists := orderBook.GetOrders()[sellOrderID]
		require.False(t, sellOrderExists, "Order ID %d should not exist in the Orders map", sellOrderID)
		_, buyOrderExists := orderBook.GetOrders()[buyOrderID]
		require.False(t, buyOrderExists, "Order ID %d should not exist in the Orders map", buyOrderID)

		// Check if the sell order is removed from the CustomerOrders map
		require.Equal(t, 0, len(orderBook.GetCustomerOrders()[sellCustomerID]), "Expected 0 order for customer ID %d", sellCustomerID)

		// Check if the buy order is not added to the CustomerOrders map
		require.Equal(t, 0, len(orderBook.GetCustomerOrders()[buyCustomerID]), "Expected 0 order for customer ID %d", buyCustomerID)
	})

	t.Run("Submit Buy Order with Higher Price and Match Sell Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		// Prepare sell order
		sellCustomerID := uint(1995)
		sellOrderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(sellCustomerID, 90, constant.SellOrder, util.CreateGTT(1))

		// Prepare buy order should match
		buyCustomerID := uint(4953)
		buyOrderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(buyCustomerID, 91, constant.BuyOrder, util.CreateGTT(1))

		// Check if the buy order is matched and not in the heap
		require.Equal(t, 0, orderBook.GetBuyOrders().Len(), "Expected 0 buy orders in the heap")

		// Check if the sell order is removed from the heap
		require.Equal(t, 0, orderBook.GetSellOrders().Len(), "Expected 0 sell orders in the heap")

		// Check if the matched orders are removed from the Orders map
		_, sellOrderExists := orderBook.GetOrders()[sellOrderID]
		require.False(t, sellOrderExists, "Order ID %d should not exist in the Orders map", sellOrderID)
		_, buyOrderExists := orderBook.GetOrders()[buyOrderID]
		require.False(t, buyOrderExists, "Order ID %d should not exist in the Orders map", buyOrderID)

		// Check if the sell order is removed from the CustomerOrders map
		require.Equal(t, 0, len(orderBook.GetCustomerOrders()[sellCustomerID]), "Expected 0 order for customer ID %d", sellCustomerID)

		// Check if the buy order is not added to the CustomerOrders map
		require.Equal(t, 0, len(orderBook.GetCustomerOrders()[buyCustomerID]), "Expected 0 order for customer ID %d", buyCustomerID)
	})

	t.Run("Submit Sell Order with Exact Price and Match Buy Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		// Prepare buy order
		buyCustomerID := uint(4953)
		buyOrderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(buyCustomerID, 90, constant.BuyOrder, util.CreateGTT(1))

		// Prepare sell order that should match
		sellCustomerID := uint(1995)
		sellOrderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(sellCustomerID, 90, constant.SellOrder, util.CreateGTT(1))

		// Check if the sell order is matched and not in the heap
		require.Equal(t, 0, orderBook.GetSellOrders().Len(), "Expected 0 sell orders in the heap")

		// Check if the buy order is removed from the heap
		require.Equal(t, 0, orderBook.GetBuyOrders().Len(), "Expected 0 buy orders in the heap")

		// Check if the matched orders are removed from the Orders map
		_, sellOrderExists := orderBook.GetOrders()[sellOrderID]
		require.False(t, sellOrderExists, "Order ID %d should not exist in the Orders map", sellOrderID)
		_, buyOrderExists := orderBook.GetOrders()[buyOrderID]
		require.False(t, buyOrderExists, "Order ID %d should not exist in the Orders map", buyOrderID)

		// Check if the sell order is removed from the CustomerOrders map
		require.Equal(t, 0, len(orderBook.GetCustomerOrders()[sellCustomerID]), "Expected 0 orders for customer ID %d", sellCustomerID)

		// Check if the buy order is removed from the CustomerOrders map
		require.Equal(t, 0, len(orderBook.GetCustomerOrders()[buyCustomerID]), "Expected 0 orders for customer ID %d", buyCustomerID)
	})

	t.Run("Submit Sell Order with Lower Price and Match Buy Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		// Prepare buy order
		buyCustomerID := uint(4953)
		buyOrderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(buyCustomerID, 90, constant.BuyOrder, util.CreateGTT(1))

		// Prepare sell order that should match
		sellCustomerID := uint(1995)
		sellOrderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(sellCustomerID, 89, constant.SellOrder, util.CreateGTT(1))

		// Check if the sell order is matched and not in the heap
		require.Equal(t, 0, orderBook.GetSellOrders().Len(), "Expected 0 sell orders in the heap")

		// Check if the buy order is removed from the heap
		require.Equal(t, 0, orderBook.GetBuyOrders().Len(), "Expected 0 buy orders in the heap")

		// Check if the matched orders are removed from the Orders map
		_, sellOrderExists := orderBook.GetOrders()[sellOrderID]
		require.False(t, sellOrderExists, "Order ID %d should not exist in the Orders map", sellOrderID)
		_, buyOrderExists := orderBook.GetOrders()[buyOrderID]
		require.False(t, buyOrderExists, "Order ID %d should not exist in the Orders map", buyOrderID)

		// Check if the sell order is removed from the CustomerOrders map
		require.Equal(t, 0, len(orderBook.GetCustomerOrders()[sellCustomerID]), "Expected 0 orders for customer ID %d", sellCustomerID)

		// Check if the buy order is removed from the CustomerOrders map
		require.Equal(t, 0, len(orderBook.GetCustomerOrders()[buyCustomerID]), "Expected 0 orders for customer ID %d", buyCustomerID)
	})

	t.Run("Submit Order with Nil GTT", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		orderID := orderBook.GetNextOrderID()
		customerID := uint(911)
		orderBook.SubmitOrder(customerID, 110, constant.BuyOrder, nil)

		// Check if the order is added to the BuyOrders heap
		require.Equal(t, 1, orderBook.GetBuyOrders().Len(), "Expected 1 buy order in the heap")

		// Check if the order is added to the Orders map
		_, exists := orderBook.GetOrders()[orderID]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID)

		// Check if the order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.GetCustomerOrders()[customerID]), "Expected 1 order for customer ID %d", customerID)
	})

	t.Run("Submit Multiple Orders from Same Customer", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		customerID := uint(9947)
		orderBook.SubmitOrder(customerID, 120, constant.SellOrder, util.CreateGTT(2))
		orderBook.SubmitOrder(customerID, 130, constant.SellOrder, util.CreateGTT(3))

		// Check if the orders are added to the SellOrders heap
		require.Equal(t, 2, orderBook.GetSellOrders().Len(), "Expected 2 sell orders in the heap")

		// Check if the orders are added to the CustomerOrders map
		require.Equal(t, 2, len(orderBook.GetCustomerOrders()[customerID]), "Expected 2 orders for customer ID %d", customerID)
	})

	t.Run("Submit Order with Expired GTT", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		expiredGTT := time.Now().Add(-1 * time.Hour)
		orderID := orderBook.GetNextOrderID()
		customerID := uint(18111995)
		orderBook.SubmitOrder(customerID, 95, constant.SellOrder, &expiredGTT)

		// Check if the order is not added to the SellOrders heap
		require.Equal(t, 0, orderBook.GetSellOrders().Len(), "Expected 0 sell orders in the heap")

		// Check if the order is not added to the Orders map
		_, exists := orderBook.GetOrders()[orderID]
		require.False(t, exists, "Order ID %d should not exist in the Orders map", orderID)

		// Check if the order is not added to the CustomerOrders map
		require.Equal(t, 0, len(orderBook.GetCustomerOrders()[customerID]), "Expected 0 orders for customer ID %d", customerID)
	})

	t.Run("Submit Orders from Same Customer with Same Prices but Different Timestamp", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		customerID := uint(69)
		orderID1 := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(customerID, 100, constant.BuyOrder, util.CreateGTT(1))

		// Change the timestamp but same price
		orderID2 := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(customerID, 100, constant.BuyOrder, util.CreateGTT(2))

		// Check if both orders are added to the BuyOrders heap
		require.Equal(t, 2, orderBook.GetBuyOrders().Len(), "Expected 2 buy orders in the heap")

		// Check if both orders are added to the Orders map
		_, exists1 := orderBook.GetOrders()[orderID1]
		require.True(t, exists1, "Order ID %d should exist in the Orders map", orderID1)
		_, exists2 := orderBook.GetOrders()[orderID2]
		require.True(t, exists2, "Order ID %d should exist in the Orders map", orderID2)

		// Check if both orders are added to the CustomerOrders map
		require.Equal(t, 2, len(orderBook.GetCustomerOrders()[customerID]), "Expected 2 orders for customer ID %d", customerID)
	})

	t.Run("Match Cancelled Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		// Submit and then cancel a buy order
		buyCustomerID := uint(3456)
		buyOrderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(buyCustomerID, 95, constant.BuyOrder, util.CreateGTT(1))
		orderBook.CancelOrder(buyOrderID)

		// Submit a sell order with matching price
		sellCustomerID := uint(7890)
		sellOrderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(sellCustomerID, 95, constant.SellOrder, util.CreateGTT(1))

		// Check if the cancelled buy order is not matched
		require.Equal(t, 0, orderBook.GetBuyOrders().Len(), "Expected 0 buy orders in the heap after cancellation")
		require.Equal(t, 1, orderBook.GetSellOrders().Len(), "Expected 1 sell order in the heap after attempting to match with cancelled buy order")

		// Check if the sell order is still present in the Orders map
		_, sellOrderExists := orderBook.GetOrders()[sellOrderID]
		require.True(t, sellOrderExists, "Sell order ID %d should exist in the Orders map", sellOrderID)

		// Check if the cancelled buy order is not in the Orders map
		_, buyOrderExists := orderBook.GetOrders()[buyOrderID]
		require.False(t, buyOrderExists, "Cancelled buy order ID %d should not exist in the Orders map", buyOrderID)

		// Check if the sell order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.GetCustomerOrders()[sellCustomerID]), "Expected 1 order for customer ID %d", sellCustomerID)

		// Check if the cancelled buy order is not in the CustomerOrders map
		require.Equal(t, 0, len(orderBook.GetCustomerOrders()[buyCustomerID]), "Expected 0 orders for customer ID %d", buyCustomerID)
	})
}

func TestOrderBookUCase_CancelOrder(t *testing.T) {
	logger := util.SetupLogger()
	defer logger.Sync() // Flushes buffer, if any
	t.Run("Cancel existing order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		orderID := orderBook.GetNextOrderID()

		// Submit initial order
		customerID := uint(123)
		orderBook.SubmitOrder(customerID, 100, constant.BuyOrder, util.CreateGTT(1))

		// Check if the order is added to the Orders map
		_, exists := orderBook.GetOrders()[orderID]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID)

		// Check if the order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.GetCustomerOrders()[customerID]), "Unexpected number of orders for customer ID %d", customerID)

		// Perform cancel existing order
		err := orderBook.CancelOrder(orderID)
		require.NoError(t, err, "CancelOrder should not return an error")

		// Check if the order is removed from the Orders map
		_, exists = orderBook.GetOrders()[orderID]
		require.False(t, exists, "Order ID %d should not exist in the Orders map", orderID)

		// Check if the order is removed from the CustomerOrders map
		require.Equal(t, 0, len(orderBook.GetCustomerOrders()[customerID]), "Unexpected number of orders for customer ID %d", customerID)
	})
	t.Run("Cancel non-existent order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		orderID := orderBook.GetNextOrderID()

		// Submit initial order
		customerID := uint(456)
		orderBook.SubmitOrder(customerID, 100, constant.BuyOrder, util.CreateGTT(1))

		// Check if the order is added to the Orders map
		_, exists := orderBook.GetOrders()[orderID]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID)

		// Check if the order is added to the CustomerOrders map
		require.Equal(t, 1, len(orderBook.GetCustomerOrders()[customerID]), "Unexpected number of orders for customer ID %d", customerID)

		// Try to cancel a non-existent order
		nonExistentOrderID := orderID + 1
		err := orderBook.CancelOrder(nonExistentOrderID)
		require.Error(t, err, fmt.Sprintf("order not found: %d", nonExistentOrderID))

		// Check if the order is still present in the Orders map
		_, exists = orderBook.GetOrders()[orderID]
		require.True(t, exists, "Order ID %d should still exist in the Orders map", orderID)

		// Check if the order is still present in the CustomerOrders map
		require.Equal(t, 1, len(orderBook.GetCustomerOrders()[customerID]), "Unexpected number of orders for customer ID %d", customerID)
	})

	t.Run("Cancel order in multi-order customer", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		// Submit initial orders
		customerID := uint(789)
		orderID1 := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(customerID, 150, constant.BuyOrder, util.CreateGTT(1))
		orderID2 := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(customerID, 160, constant.BuyOrder, util.CreateGTT(1))

		// Cancel the first order
		err := orderBook.CancelOrder(orderID1)
		require.NoError(t, err, "CancelOrder should not return an error")

		// Check if the first order is removed from the Orders map
		_, exists := orderBook.GetOrders()[orderID1]
		require.False(t, exists, "Order ID %d should not exist in the Orders map", orderID1)

		// Check if the second order is still present in the Orders map
		_, exists = orderBook.GetOrders()[orderID2]
		require.True(t, exists, "Order ID %d should exist in the Orders map", orderID2)

		// Check if the CustomerOrders map is updated correctly
		require.Equal(t, 1, len(orderBook.GetCustomerOrders()[customerID]), "Unexpected number of orders for customer ID %d", customerID)
	})
}

func TestOrderBookUCase_QueryOrders(t *testing.T) {
	logger := util.SetupLogger()
	defer logger.Sync() // Flushes buffer, if any
	t.Run("Query Orders with Active Orders", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		customerID := uint(100)
		orderBook.SubmitOrder(customerID, 120, constant.BuyOrder, util.CreateGTT(1))
		orderBook.SubmitOrder(customerID, 130, constant.BuyOrder, util.CreateGTT(2))

		// Query active orders
		orders := orderBook.QueryOrders(customerID)

		require.Equal(t, 2, len(orders), "Expected 2 active orders for customer ID %d", customerID)
	})

	t.Run("Query Orders with Expired Orders", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		customerID := uint(101)
		expiredGTT := time.Now().Add(-1 * time.Hour)
		orderBook.SubmitOrder(customerID, 120, constant.BuyOrder, &expiredGTT)
		orderBook.SubmitOrder(customerID, 130, constant.BuyOrder, &expiredGTT)

		// Query orders
		orders := orderBook.QueryOrders(customerID)

		require.Equal(t, 0, len(orders), "Expected 0 active order for customer ID %d", customerID)
	})

	t.Run("Query Orders with No Orders", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		customerID := uint(102)

		// Query orders
		orders := orderBook.QueryOrders(customerID)

		require.Equal(t, 0, len(orders), "Expected 0 orders for customer ID %d", customerID)
	})

	t.Run("Query Orders with Both Active and Expired Orders", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		customerID := uint(103)
		expiredGTT := time.Now().Add(-1 * time.Hour)
		orderBook.SubmitOrder(customerID, 140, constant.BuyOrder, &expiredGTT)
		activeOrderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(customerID, 150, constant.BuyOrder, util.CreateGTT(1))

		// Query orders
		orders := orderBook.QueryOrders(customerID)

		require.Equal(t, 1, len(orders), "Expected 1 active order for customer ID %d", customerID)
		require.Equal(t, activeOrderID, orders[0].ID, "Expected order with ID %d", activeOrderID)
	})

	t.Run("Query Orders After Canceling an Order", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		customerID := uint(104)
		orderID := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(customerID, 150, constant.BuyOrder, util.CreateGTT(1))
		orderBook.CancelOrder(orderID)

		// Query orders
		orders := orderBook.QueryOrders(customerID)

		require.Equal(t, 0, len(orders), "Expected 0 orders for customer ID %d after cancellation", customerID)
	})

	t.Run("Query Orders After Canceling All Orders", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		customerID := uint(105)
		orderID1 := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(customerID, 160, constant.BuyOrder, util.CreateGTT(1))
		orderID2 := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(customerID, 170, constant.BuyOrder, util.CreateGTT(2))
		orderBook.CancelOrder(orderID1)
		orderBook.CancelOrder(orderID2)

		// Query orders
		orders := orderBook.QueryOrders(customerID)

		require.Equal(t, 0, len(orders), "Expected 0 orders for customer ID %d after canceling all orders", customerID)
	})

	t.Run("Query Orders with Cancelled Orders but Active Orders Present", func(t *testing.T) {
		// Create a new order book
		orderBook := module.NewOrderBookUCase(logger)

		customerID := uint(106)
		expiredGTT := time.Now().Add(-1 * time.Hour)
		orderID1 := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(customerID, 180, constant.BuyOrder, &expiredGTT)
		orderID2 := orderBook.GetNextOrderID()
		orderBook.SubmitOrder(customerID, 190, constant.BuyOrder, util.CreateGTT(1))
		orderBook.CancelOrder(orderID1)

		// Query orders
		orders := orderBook.QueryOrders(customerID)

		require.Equal(t, 1, len(orders), "Expected 1 active order for customer ID %d", customerID)
		require.Equal(t, orderID2, orders[0].ID, "Expected order with ID %d", orderID2)
	})
}
