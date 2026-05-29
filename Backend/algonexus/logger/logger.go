package logger

import (
	"os"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

		// Level overridable via LOG_LEVEL (e.g. "error" to silence hot-path INFO logs
		// during load tests). Defaults to production INFO.
		if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
			if parsed, perr := zapcore.ParseLevel(lvl); perr == nil {
				config.Level = zap.NewAtomicLevelAt(parsed)
			}
		}

		logger, _ := config.Build()
		instance = &Logger{zapLogger: logger}
	})
	return instance
}

func CallerInfo(fields ...zap.Field) []zap.Field {
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		fields = append(fields,
			zap.String("caller_func", funcName),
			zap.String("caller_file", file),
			zap.Int("caller_line", line),
		)
		return fields
	}
	return nil
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	fields = CallerInfo(fields...)
	fields = append(fields, zap.Time("logged_at", time.Now()))
	l.zapLogger.Debug(msg, fields...)
}

// Info logs an informational message with timestamp
func (l *Logger) Info(msg string, fields ...zap.Field) {
	fields = append(fields, zap.Time("logged_at", time.Now()))

	l.zapLogger.Info(msg, fields...)
}

// Warning logs an informational message with timestamp
func (l *Logger) Warning(msg string, fields ...zap.Field) {
	fields = append(fields, zap.Time("logged_at", time.Now()))
	l.zapLogger.Warn(msg, fields...)
}

// Error logs an error message with timestamp
func (l *Logger) Error(msg string, fields ...zap.Field) {
	fields = append(fields, zap.Time("logged_at", time.Now()))
	l.zapLogger.Error(msg, fields...)
}

// Error logs an fatal message with timestamp
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	fields = append(fields, zap.Time("logged_at", time.Now()))
	l.zapLogger.Fatal(msg, fields...)
}
