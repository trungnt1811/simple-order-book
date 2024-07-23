package model

// OrderHeap is a priority queue for orders, implemented as a container/heap.
type OrderHeap struct {
	Orders []*Order
	Desc   bool // Descending order (for buy orders)
}

// Len returns the number of orders in the heap.
func (h OrderHeap) Len() int {
	return len(h.Orders)
}

// Less compares two orders in the heap.
func (h OrderHeap) Less(i, j int) bool {
	if h.Desc {
		return h.Orders[i].Price > h.Orders[j].Price
	}
	return h.Orders[i].Price < h.Orders[j].Price
}

// Swap swaps two orders in the heap.
func (h OrderHeap) Swap(i, j int) {
	h.Orders[i], h.Orders[j] = h.Orders[j], h.Orders[i]
}

// Push adds an order to the heap.
func (h *OrderHeap) Push(x interface{}) { h.Orders = append(h.Orders, x.(*Order)) }

// Pop removes and returns the highest priority order from the heap.
func (h *OrderHeap) Pop() interface{} {
	old := h.Orders
	n := len(old)
	order := old[n-1]
	h.Orders = old[0 : n-1]
	return order
}
