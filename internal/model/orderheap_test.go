package model_test

import (
	"container/heap"
	"testing"

	"github.com/trungnt1811/simple-order-book/internal/model"
)

// TestOrderHeap tests the OrderHeap implementation.
func TestOrderHeap(t *testing.T) {
	testCases := []struct {
		name          string
		desc          bool
		orders        []*model.Order
		expectedOrder []int
	}{
		{
			name: "Ascending order test (sell orders)",
			desc: false,
			orders: []*model.Order{
				{Price: 100},
				{Price: 99},
				{Price: 101},
				{Price: 98},
			},
			expectedOrder: []int{98, 99, 100, 101},
		},
		{
			name: "Descending order test (buy orders)",
			desc: true,
			orders: []*model.Order{
				{Price: 100},
				{Price: 99},
				{Price: 101},
				{Price: 98},
			},
			expectedOrder: []int{101, 100, 99, 98},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			orderHeap := &model.OrderHeap{Desc: tc.desc}
			for _, order := range tc.orders {
				heap.Push(orderHeap, order)
			}

			for i, expectedPrice := range tc.expectedOrder {
				order := heap.Pop(orderHeap).(*model.Order)
				if order.Price != expectedPrice {
					t.Errorf(
						"Test case failed: %s. Order %d: expected price %d, got %d",
						tc.name,
						i+1,
						expectedPrice,
						order.Price,
					)
				}
			}
		})
	}
}
