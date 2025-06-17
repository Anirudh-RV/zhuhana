package middleware

import (
	"time"

	"algonexus/logger" // Replace with the actual module path

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger logs details about each HTTP request including start time
func RequestLogger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// After request
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		log.Info("API Request",
			zap.Time("start_time", start),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.Int("status", statusCode),
			zap.Duration("duration", duration),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}
