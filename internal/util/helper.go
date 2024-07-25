package util

import (
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/trungnt1811/simple-order-book/internal/logger"
)

// Helper function to create a GTT time.
func CreateGTT(hours int) *time.Time {
	gtt := time.Now().Add(time.Duration(hours) * time.Hour)
	return &gtt
}

// Helper function to setup logger.
func SetupLogger() *zap.Logger {
	logger, err := logger.Setup()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}
	return logger
}
