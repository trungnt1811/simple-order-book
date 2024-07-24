package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Configure initializes and configures a zap logger.
func Setup() (*zap.Logger, error) {
	// Default level is info
	logLevel := zapcore.InfoLevel

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(logLevel)

	return config.Build()
}
