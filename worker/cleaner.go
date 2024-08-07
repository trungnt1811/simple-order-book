package worker

import (
	"time"

	"github.com/trungnt1811/simple-order-book/internal/interfaces"
)

// cleaner is responsible for periodically removing expired orders
// from the order book.
type cleaner struct {
	OrderBook interfaces.OrderBookUCase
}

// NewCleaner creates a new cleaner instance with the provided
// order book use case.
func NewCleaner(orderBook interfaces.OrderBookUCase) cleaner {
	return cleaner{
		OrderBook: orderBook,
	}
}

// RemoveExpiredBuyOrders starts a ticker that triggers the removal
// of expired buy orders every 5 seconds.
func (c *cleaner) RemoveExpiredBuyOrders() {
	// Create a new ticker that ticks every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop() // Ensure the ticker is stopped when the function exits

	// Iterate over the ticker's channel to trigger the removal
	for range ticker.C {
		c.OrderBook.RemoveExpiredBuyOrders()
	}

	select {} // Block forever to keep the goroutine running
}

// RemoveExpiredSellOrders starts a ticker that triggers the removal
// of expired sell orders every 5 seconds.
func (c *cleaner) RemoveExpiredSellOrders() {
	// Create a new ticker that ticks every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop() // Ensure the ticker is stopped when the function exits

	// Iterate over the ticker's channel to trigger the removal
	for range ticker.C {
		c.OrderBook.RemoveExpiredSellOrders()
	}

	select {} // Block forever to keep the goroutine running
}
