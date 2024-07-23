package model_test

import (
	"container/heap"
	"testing"

	"github.com/trungnt1811/simple-order-book/internal/model"
)

type ExpectedOrder struct {
	CustomerID int
	Price      int
}

// TestOrderHeap tests the OrderHeap implementation.
func TestOrderHeap(t *testing.T) {
	testCases := []struct {
		name           string
		desc           bool
		orders         []*model.Order
		expectedOrders []ExpectedOrder
	}{
		{
			name: "Ascending order test (sell orders)",
			desc: false,
			orders: []*model.Order{
				{CustomerID: 1, Price: 100},
				{CustomerID: 2, Price: 99},
				{CustomerID: 3, Price: 101},
				{CustomerID: 4, Price: 98},
			},
			expectedOrders: []ExpectedOrder{
				{
					CustomerID: 4,
					Price:      98,
				},
				{
					CustomerID: 2,
					Price:      99,
				},
				{
					CustomerID: 1,
					Price:      100,
				},
				{
					CustomerID: 3,
					Price:      101,
				},
			},
		},
		{
			name: "Descending order test (buy orders)",
			desc: true,
			orders: []*model.Order{
				{CustomerID: 1, Price: 100},
				{CustomerID: 2, Price: 99},
				{CustomerID: 3, Price: 101},
				{CustomerID: 4, Price: 98},
			},
			expectedOrders: []ExpectedOrder{
				{
					CustomerID: 3,
					Price:      101,
				},
				{
					CustomerID: 1,
					Price:      100,
				},
				{
					CustomerID: 2,
					Price:      99,
				},
				{
					CustomerID: 4,
					Price:      98,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			orderHeap := &model.OrderHeap{Desc: tc.desc}
			for _, order := range tc.orders {
				heap.Push(orderHeap, order)
			}

			for i, expectedOrder := range tc.expectedOrders {
				order := heap.Pop(orderHeap).(*model.Order)
				if order.CustomerID != expectedOrder.CustomerID {
					t.Errorf(
						"Test case failed: %s. Order %d: expected customer ID %d, got %d",
						tc.name,
						i+1,
						expectedOrder.CustomerID,
						order.CustomerID,
					)
				}
				if order.Price != expectedOrder.Price {
					t.Errorf(
						"Test case failed: %s. Order %d: expected price %d, got %d",
						tc.name,
						i+1,
						expectedOrder.Price,
						order.Price,
					)
				}
			}
		})
	}
}
