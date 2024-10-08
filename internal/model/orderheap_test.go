package model_test

import (
	"container/heap"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trungnt1811/simple-order-book/internal/constant"
	"github.com/trungnt1811/simple-order-book/internal/model"
)

type expectedOrder struct {
	CustomerID uint
	Price      uint
}

// TestOrderHeap tests the OrderHeap implementation.
func TestOrderHeap(t *testing.T) {
	testCases := []struct {
		name           string
		orderType      constant.OrderType
		orders         []*model.Order
		expectedOrders []expectedOrder
	}{
		{
			name:      "Sell orders test",
			orderType: constant.SellOrder,
			orders: []*model.Order{
				{CustomerID: 1, Price: 100},
				{CustomerID: 2, Price: 99},
				{CustomerID: 3, Price: 101},
				{CustomerID: 4, Price: 98},
			},
			expectedOrders: []expectedOrder{
				{CustomerID: 4, Price: 98},
				{CustomerID: 2, Price: 99},
				{CustomerID: 1, Price: 100},
				{CustomerID: 3, Price: 101},
			},
		},
		{
			name:      "Buy orders test",
			orderType: constant.BuyOrder,
			orders: []*model.Order{
				{CustomerID: 1, Price: 100},
				{CustomerID: 2, Price: 99},
				{CustomerID: 3, Price: 101},
				{CustomerID: 4, Price: 98},
			},
			expectedOrders: []expectedOrder{
				{CustomerID: 3, Price: 101},
				{CustomerID: 1, Price: 100},
				{CustomerID: 2, Price: 99},
				{CustomerID: 4, Price: 98},
			},
		},
		{
			name:           "Empty heap test",
			orderType:      constant.BuyOrder,
			orders:         []*model.Order{},
			expectedOrders: []expectedOrder{},
		},
		{
			name:      "Duplicate prices with various timestamps test (sell orders)",
			orderType: constant.SellOrder,
			orders: []*model.Order{
				{CustomerID: 1, Price: 100, Timestamp: time.Now()},
				{CustomerID: 2, Price: 100, Timestamp: time.Now().Add(1 * time.Second)},
				{CustomerID: 3, Price: 100, Timestamp: time.Now().Add(2 * time.Second)},
				{CustomerID: 2, Price: 99, Timestamp: time.Now().Add(3 * time.Second)},
				{CustomerID: 4, Price: 99, Timestamp: time.Now().Add(4 * time.Second)},
				{CustomerID: 5, Price: 101, Timestamp: time.Now().Add(5 * time.Second)},
			},
			expectedOrders: []expectedOrder{
				{CustomerID: 2, Price: 99},
				{CustomerID: 4, Price: 99},
				{CustomerID: 1, Price: 100},
				{CustomerID: 2, Price: 100},
				{CustomerID: 3, Price: 100},
				{CustomerID: 5, Price: 101},
			},
		},
		{
			name:      "Duplicate prices with various timestamps test (buy orders)",
			orderType: constant.BuyOrder,
			orders: []*model.Order{
				{CustomerID: 1, Price: 100, Timestamp: time.Now()},
				{CustomerID: 2, Price: 100, Timestamp: time.Now().Add(1 * time.Second)},
				{CustomerID: 3, Price: 100, Timestamp: time.Now().Add(2 * time.Second)},
				{CustomerID: 2, Price: 99, Timestamp: time.Now().Add(3 * time.Second)},
				{CustomerID: 4, Price: 99, Timestamp: time.Now().Add(4 * time.Second)},
				{CustomerID: 5, Price: 101, Timestamp: time.Now().Add(5 * time.Second)},
			},
			expectedOrders: []expectedOrder{
				{CustomerID: 5, Price: 101},
				{CustomerID: 1, Price: 100},
				{CustomerID: 2, Price: 100},
				{CustomerID: 3, Price: 100},
				{CustomerID: 2, Price: 99},
				{CustomerID: 4, Price: 99},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			orderHeap := &model.OrderHeap{Type: tc.orderType}
			for _, order := range tc.orders {
				heap.Push(orderHeap, order)
			}

			for i, expectedOrder := range tc.expectedOrders {
				order := heap.Pop(orderHeap).(*model.Order)
				require.Equal(t, expectedOrder.CustomerID, order.CustomerID, "[%s] Order %d: expected customer ID %d, got %d", tc.name, i+1, expectedOrder.CustomerID, order.CustomerID)
				require.Equal(t, expectedOrder.Price, order.Price, "[%s] Order %d: expected price %d, got %d", tc.name, i+1, expectedOrder.Price, order.Price)
			}
		})
	}
}

// TestOrderHeapPopAndPush tests the OrderHeap implementation for popping and pushing orders.
func TestOrderHeap_PopAndPush(t *testing.T) {
	t.Run("Push to empty heap", func(t *testing.T) {
		orderHeap := &model.OrderHeap{Type: constant.SellOrder}
		order := &model.Order{
			CustomerID: 1,
			Price:      100,
			Timestamp:  time.Now(),
		}
		heap.Push(orderHeap, order)

		require.Equal(t, 1, orderHeap.Len(), "Expected heap length 1, got %d", orderHeap.Len())

		poppedOrder := heap.Pop(orderHeap).(*model.Order)
		require.Equal(t, order.CustomerID, poppedOrder.CustomerID, "Expected order with CustomerID %d and Price %d, got CustomerID %d and Price %d", order.CustomerID, order.Price, poppedOrder.CustomerID, poppedOrder.Price)
		require.Equal(t, order.Price, poppedOrder.Price, "Expected order with CustomerID %d and Price %d, got CustomerID %d and Price %d", order.CustomerID, order.Price, poppedOrder.CustomerID, poppedOrder.Price)
	})

	t.Run("Push to non-empty heap (sell orders)", func(t *testing.T) {
		orderHeap := &model.OrderHeap{Type: constant.SellOrder}
		orders := []*model.Order{
			{CustomerID: 1, Price: 100, Timestamp: time.Now()},
			{CustomerID: 2, Price: 99, Timestamp: time.Now()},
		}
		for _, order := range orders {
			heap.Push(orderHeap, order)
		}
		pushOrder := &model.Order{
			CustomerID: 3,
			Price:      98,
			Timestamp:  time.Now(),
		}
		heap.Push(orderHeap, pushOrder)

		expectedOrders := []expectedOrder{
			{CustomerID: 3, Price: 98},
			{CustomerID: 2, Price: 99},
			{CustomerID: 1, Price: 100},
		}

		for i, expectedOrder := range expectedOrders {
			order := heap.Pop(orderHeap).(*model.Order)
			require.Equal(t, expectedOrder.CustomerID, order.CustomerID, "Order %d: expected CustomerID %d and Price %d, got CustomerID %d and Price %d", i+1, expectedOrder.CustomerID, expectedOrder.Price, order.CustomerID, order.Price)
			require.Equal(t, expectedOrder.Price, order.Price, "Order %d: expected CustomerID %d and Price %d, got CustomerID %d and Price %d", i+1, expectedOrder.CustomerID, expectedOrder.Price, order.CustomerID, order.Price)
		}
	})

	t.Run("Push to non-empty heap (buy orders)", func(t *testing.T) {
		orderHeap := &model.OrderHeap{Type: constant.BuyOrder}
		orders := []*model.Order{
			{CustomerID: 1, Price: 100, Timestamp: time.Now()},
			{CustomerID: 2, Price: 99, Timestamp: time.Now()},
		}
		for _, order := range orders {
			heap.Push(orderHeap, order)
		}
		pushOrder := &model.Order{
			CustomerID: 3,
			Price:      101,
			Timestamp:  time.Now(),
		}
		heap.Push(orderHeap, pushOrder)

		expectedOrders := []expectedOrder{
			{CustomerID: 3, Price: 101},
			{CustomerID: 1, Price: 100},
			{CustomerID: 2, Price: 99},
		}

		for i, expectedOrder := range expectedOrders {
			order := heap.Pop(orderHeap).(*model.Order)
			require.Equal(t, expectedOrder.CustomerID, order.CustomerID, "Order %d: expected CustomerID %d and Price %d, got CustomerID %d and Price %d", i+1, expectedOrder.CustomerID, expectedOrder.Price, order.CustomerID, order.Price)
			require.Equal(t, expectedOrder.Price, order.Price, "Order %d: expected CustomerID %d and Price %d, got CustomerID %d and Price %d", i+1, expectedOrder.CustomerID, expectedOrder.Price, order.CustomerID, order.Price)
		}
	})

	t.Run("Pop from heap and push new order (sell orders)", func(t *testing.T) {
		orderHeap := &model.OrderHeap{Type: constant.SellOrder}
		orders := []*model.Order{
			{CustomerID: 1, Price: 100, Timestamp: time.Now()},
			{CustomerID: 2, Price: 99, Timestamp: time.Now()},
		}
		for _, order := range orders {
			heap.Push(orderHeap, order)
		}
		heap.Pop(orderHeap)
		pushOrder := &model.Order{
			CustomerID: 3,
			Price:      101,
			Timestamp:  time.Now(),
		}
		heap.Push(orderHeap, pushOrder)

		expectedOrders := []expectedOrder{
			{CustomerID: 1, Price: 100},
			{CustomerID: 3, Price: 101},
		}

		for i, expectedOrder := range expectedOrders {
			order := heap.Pop(orderHeap).(*model.Order)
			require.Equal(t, expectedOrder.CustomerID, order.CustomerID, "Order %d: expected CustomerID %d and Price %d, got CustomerID %d and Price %d", i+1, expectedOrder.CustomerID, expectedOrder.Price, order.CustomerID, order.Price)
			require.Equal(t, expectedOrder.Price, order.Price, "Order %d: expected CustomerID %d and Price %d, got CustomerID %d and Price %d", i+1, expectedOrder.CustomerID, expectedOrder.Price, order.CustomerID, order.Price)
		}
	})

	t.Run("Pop from heap and push new order (buy orders)", func(t *testing.T) {
		orderHeap := &model.OrderHeap{Type: constant.BuyOrder}
		orders := []*model.Order{
			{CustomerID: 1, Price: 100, Timestamp: time.Now()},
			{CustomerID: 2, Price: 99, Timestamp: time.Now()},
		}
		for _, order := range orders {
			heap.Push(orderHeap, order)
		}
		heap.Pop(orderHeap)
		pushOrder := &model.Order{
			CustomerID: 3,
			Price:      98,
			Timestamp:  time.Now(),
		}
		heap.Push(orderHeap, pushOrder)

		expectedOrders := []expectedOrder{
			{CustomerID: 2, Price: 99},
			{CustomerID: 3, Price: 98},
		}

		for i, expectedOrder := range expectedOrders {
			order := heap.Pop(orderHeap).(*model.Order)
			require.Equal(t, expectedOrder.CustomerID, order.CustomerID, "Order %d: expected CustomerID %d and Price %d, got CustomerID %d and Price %d", i+1, expectedOrder.CustomerID, expectedOrder.Price, order.CustomerID, order.Price)
			require.Equal(t, expectedOrder.Price, order.Price, "Order %d: expected CustomerID %d and Price %d, got CustomerID %d and Price %d", i+1, expectedOrder.CustomerID, expectedOrder.Price, order.CustomerID, order.Price)
		}
	})

	t.Run("Pop from heap and push order that was recently popped (sell orders)", func(t *testing.T) {
		orderHeap := &model.OrderHeap{Type: constant.SellOrder}
		orders := []*model.Order{
			{CustomerID: 1, Price: 100, Timestamp: time.Now()},
			{CustomerID: 2, Price: 99, Timestamp: time.Now()},
		}
		for _, order := range orders {
			heap.Push(orderHeap, order)
		}
		poppedOrder := heap.Pop(orderHeap).(*model.Order)
		heap.Push(orderHeap, poppedOrder)

		expectedOrders := []expectedOrder{
			{CustomerID: 2, Price: 99},
			{CustomerID: 1, Price: 100},
		}

		for i, expectedOrder := range expectedOrders {
			order := heap.Pop(orderHeap).(*model.Order)
			require.Equal(t, expectedOrder.CustomerID, order.CustomerID, "Order %d: expected CustomerID %d and Price %d, got CustomerID %d and Price %d", i+1, expectedOrder.CustomerID, expectedOrder.Price, order.CustomerID, order.Price)
			require.Equal(t, expectedOrder.Price, order.Price, "Order %d: expected CustomerID %d and Price %d, got CustomerID %d and Price %d", i+1, expectedOrder.CustomerID, expectedOrder.Price, order.CustomerID, order.Price)
		}
	})

	t.Run("Pop from heap and push order that was recently popped (buy orders)", func(t *testing.T) {
		orderHeap := &model.OrderHeap{Type: constant.BuyOrder}
		orders := []*model.Order{
			{CustomerID: 1, Price: 100, Timestamp: time.Now()},
			{CustomerID: 2, Price: 99, Timestamp: time.Now()},
		}
		for _, order := range orders {
			heap.Push(orderHeap, order)
		}
		poppedOrder := heap.Pop(orderHeap).(*model.Order)
		heap.Push(orderHeap, poppedOrder)

		expectedOrders := []expectedOrder{
			{CustomerID: 1, Price: 100},
			{CustomerID: 2, Price: 99},
		}

		for i, expectedOrder := range expectedOrders {
			order := heap.Pop(orderHeap).(*model.Order)
			require.Equal(t, expectedOrder.CustomerID, order.CustomerID, "Order %d: expected CustomerID %d and Price %d, got CustomerID %d and Price %d", i+1, expectedOrder.CustomerID, expectedOrder.Price, order.CustomerID, order.Price)
			require.Equal(t, expectedOrder.Price, order.Price, "Order %d: expected CustomerID %d and Price %d, got CustomerID %d and Price %d", i+1, expectedOrder.CustomerID, expectedOrder.Price, order.CustomerID, order.Price)
		}
	})
}
