// Logger for generating strategy statistic

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	instance *Logger
	once     sync.Once
)

// Logger struct that holds the zap logger
type Logger struct {
	zapLogger *zap.Logger
}

// NewLogger initializes and returns a singleton Logger instance
func NewLogger() *Logger {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"                   // Explicitly set timestamp key
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Format as ISO8601

		logger, _ := config.Build()
		instance = &Logger{zapLogger: logger}
	})
	return instance
}
