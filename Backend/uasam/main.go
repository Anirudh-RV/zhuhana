package main

import (
	"context"
	"uasam/cache"
	"uasam/commonutils"
	"uasam/db"
	"uasam/email"
	"uasam/logger"
	"uasam/routes"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	var ctx = context.Background()

	log := logger.NewLogger()
	go log.Info("Logger started", zap.String("Execution Level", "Root"))

	db.InitDB(log)
	go log.Info("DB connection successful", zap.String("Execution Level", "Root"))

	emailService := email.NewEmailService(&ctx, log)
	go log.Info("DB connection successful", zap.String("Execution Level", "Root"))

	jwtService := commonutils.NewJWTService(log)
	go log.Info("DB connection successful", zap.String("Execution Level", "Root"))

	cache.InitRedis(&ctx, log)
	go log.Info("Redis connection successful", zap.String("Execution Level", "Root"))

	router := gin.Default()
	go log.Info("Router setup successful", zap.String("Execution Level", "Root"))

	routes.RegisterRoutes(&ctx, router, log, db.DB, cache.RedisObj, emailService, jwtService)

	go log.Info("Starting application at port 8080...", zap.String("Execution Level", "Root"))
	router.Run(":8080")
}
