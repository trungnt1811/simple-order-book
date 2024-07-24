package model

import "github.com/trungnt1811/simple-order-book/internal/constant"

// OrderHeap is a priority queue for orders, implemented as a container/heap.
type OrderHeap struct {
	Orders []*Order
	Type   constant.OrderType // Max heap for buy order, min heap for sell order
}

// Len returns the number of orders in the heap.
func (h OrderHeap) Len() int {
	return len(h.Orders)
}

// Less compares two orders in the heap.
func (h OrderHeap) Less(i, j int) bool {
	if h.Orders[i].Price == h.Orders[j].Price {
		return h.Orders[i].Timestamp.Before(h.Orders[j].Timestamp)
	}
	if h.Type == constant.BuyOrder {
		return h.Orders[i].Price > h.Orders[j].Price
	}
	return h.Orders[i].Price < h.Orders[j].Price
}

// Swap swaps two orders in the heap.
func (h OrderHeap) Swap(i, j int) {
	h.Orders[i], h.Orders[j] = h.Orders[j], h.Orders[i]
}

// Push adds an order to the heap.
func (h *OrderHeap) Push(x interface{}) {
	h.Orders = append(h.Orders, x.(*Order))
}

// Pop removes and returns the highest priority order from the heap.
func (h *OrderHeap) Pop() interface{} {
	n := len(h.Orders)
	order := h.Orders[n-1]       // Retrieves the last element of the slice, which is the order to be removed
	h.Orders = h.Orders[0 : n-1] // Removing the last element
	return order
}
