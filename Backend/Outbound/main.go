package main

import (
	"context"
	"outbound/cache"
	"outbound/logger"
	"outbound/routes"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	var ctx = context.Background()

	router := gin.Default()
	log := logger.NewLogger()
	go log.Info("Application started", zap.String("Execution Level", "Root"))

	cache.InitRedis(ctx, log)
	go log.Info("Redis connection successful", zap.String("Execution Level", "Root"))

	routes.RegisterRoutes(router, log, cache.RedisObj)

	router.Run(":8080")
}
