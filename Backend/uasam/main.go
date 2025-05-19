package main

import (
	"context"
	"uasam/cache"
	"uasam/commonutils"
	"uasam/db"
	"uasam/email"
	"uasam/logger"
	"uasam/middleware"
	"uasam/routes"

	constants "uasam/constants"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	var ctx = context.Background()

	log := logger.NewLogger()
	go log.Info("logger started", zap.String("execution level", "Root"))

	db.InitDB(log)
	go log.Info("db connection successful", zap.String("execution level", "Root"))

	emailService := email.NewEmailService(&ctx, log)
	go log.Info("email service setup successful", zap.String("execution level", "Root"))

	jwtService := commonutils.NewJWTService(log)
	go log.Info("jwt service setup successful", zap.String("execution level", "Root"))

	cache.InitRedis(&ctx, log)
	go log.Info("cache connection successful", zap.String("execution level", "Root"))

	router := gin.Default()
	go log.Info("router setup successful", zap.String("execution level", "Root"))

	authMiddleware := middleware.AuthMiddleware(constants.API_AUTHENTICATION_ENDPOINT)
	go log.Info("authentication middleware initialization successful", zap.String("execution level", "Root"))

	routes.RegisterRoutes(&ctx, router, log, db.DB, cache.RedisObj, emailService, jwtService, authMiddleware)

	go log.Info("starting application at port 8080...", zap.String("execution level", "Root"))
	router.Run(":8080")
}
