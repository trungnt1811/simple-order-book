package util

import (
	"log"
	"time"

	"github.com/trungnt1811/simple-order-book/internal/logger"
	"github.com/trungnt1811/simple-order-book/internal/module"
)

// Helper function to create a GTT time.
func CreateGTT(hours int) *time.Time {
	gtt := time.Now().Add(time.Duration(hours) * time.Hour)
	return &gtt
}

// Helper function to create a new order book with a logger.
func NewOrderBookWithLogger() *module.OrderBook {
	logger, err := logger.Setup()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}
	defer logger.Sync() // Flushes buffer, if any
	return module.NewOrderBook(logger)
}
