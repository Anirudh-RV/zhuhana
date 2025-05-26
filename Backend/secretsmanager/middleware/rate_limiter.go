package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"secretsmanager/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RateLimiterConfig struct {
	Source      string // "body", "query", "header"
	Param       string // e.g., "emailID"
	Limit       int    // per parameter
	Window      int    // Limit window in seconds
	EnableParam bool   // if true, apply per-param limiting
	EnableIP    bool   // if true, apply per-IP limiting
	IPLimit     int    // per IP address
	IPWindow    int    // IP limit time window in seconds
	Endpoint    string // The endpoint for the RateLimiter
}

type RateLimiterResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}

func RateLimiter(redisClient *redis.Client, log *logger.Logger, config RateLimiterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		paramWindow := time.Duration(config.Window) * time.Second
		ipWindow := time.Duration(config.IPWindow) * time.Second

		var paramValue string
		if config.EnableParam {
			var ok bool
			switch config.Source {
			case "query":
				paramValue = c.Query(config.Param)
				ok = paramValue != ""
			case "header":
				paramValue = c.GetHeader(config.Param)
				ok = paramValue != ""
			case "body":
				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, RateLimiterResponse{Status: -1, StatusDescription: "Invalid request body"})
					return
				}
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // reset body for next handler

				var bodyMap map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &bodyMap); err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, RateLimiterResponse{Status: -1, StatusDescription: "Invalid JSON body"})
					return
				}
				v, exists := bodyMap[config.Param]
				if exists {
					paramValue, ok = v.(string)
				}
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, RateLimiterResponse{Status: 0, StatusDescription: "Unsupported source"})
				return
			}
			if !ok || paramValue == "" {
				go log.Warning(fmt.Sprintf("Missing parameter: %s", config.Param), zap.String("execution level", "RateLimiter"))
				c.AbortWithStatusJSON(http.StatusBadRequest, RateLimiterResponse{Status: -1, StatusDescription: "Missing parameters"})
				return
			}
		}

		clientIP := c.ClientIP()
		pipe := redisClient.TxPipeline()

		if config.EnableParam {
			paramKey := fmt.Sprintf("rl:endpoint:%s:param:%s:%s", config.Endpoint, config.Param, paramValue)
			count, err := redisClient.Get(ctx, paramKey).Int()
			if err != nil && err != redis.Nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, RateLimiterResponse{Status: 0, StatusDescription: "Rate limiter error"})
				return
			}
			if count >= config.Limit {
				go log.Error(fmt.Sprintf("Too many requests for %s=%s", config.Param, paramValue), zap.String("execution level", "RateLimiter"))
				c.AbortWithStatusJSON(http.StatusTooManyRequests, RateLimiterResponse{Status: -1, StatusDescription: "Too many requests. Try again after sometime."})
				return
			}
			pipe.Incr(ctx, paramKey)
			if count == 0 {
				pipe.Expire(ctx, paramKey, paramWindow)
			}
		}

		if config.EnableIP {
			ipKey := fmt.Sprintf("rl:endpoint:%s:ip:%s", config.Endpoint, clientIP)
			ipCount, err := redisClient.Get(ctx, ipKey).Int()
			if err != nil && err != redis.Nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, RateLimiterResponse{Status: 0, StatusDescription: "Rate limiter error"})
				return
			}
			if ipCount >= config.IPLimit {
				go log.Error(fmt.Sprintf("Too many requests from IP: %s", clientIP), zap.String("execution level", "RateLimiter"))
				c.AbortWithStatusJSON(http.StatusTooManyRequests, RateLimiterResponse{Status: -1, StatusDescription: "Too many requests. Try again after sometime."})
				return
			}
			pipe.Incr(ctx, ipKey)
			if ipCount == 0 {
				pipe.Expire(ctx, ipKey, ipWindow)
			}
		}

		_, _ = pipe.Exec(ctx)
		c.Next()
	}
}
