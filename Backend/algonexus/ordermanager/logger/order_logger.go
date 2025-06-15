package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *OrderLogger
	once     sync.Once
)

// Logger struct that holds the zap logger
type OrderLogger struct {
	zapLogger *zap.Logger
}

// NewLogger initializes and returns a singleton Logger instance
func NewOrderLogger() *OrderLogger {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"                   // Explicitly set timestamp key
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Format as ISO8601

		logger, _ := config.Build()
		instance = &OrderLogger{zapLogger: logger}
	})
	return instance
}
