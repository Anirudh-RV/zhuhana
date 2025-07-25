package main

import (
	"context"
	"secretsmanager/cache"
	"secretsmanager/constants"
	"secretsmanager/db"
	"secretsmanager/logger"
	"secretsmanager/middleware"
	"secretsmanager/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	var ctx = context.Background()

	log := logger.NewLogger()
	go log.Info("Logger started", zap.String("Execution Level", "Root"))

	db.InitDB(log)
	go log.Info("DB connection successful", zap.String("Execution Level", "Root"))

	cache.InitRedis(ctx, log)
	go log.Info("Redis connection successful", zap.String("Execution Level", "Root"))

	router := gin.Default()
	go log.Info("Router setup successful", zap.String("Execution Level", "Root"))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "USER_TOKEN"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	go log.Info("Router CORS configuration setup successful", zap.String("Execution Level", "Root"))

	authMiddleware := middleware.AuthMiddleware(constants.API_AUTHENTICATION_ENDPOINT)
	go log.Info("authentication middleware initialization successful", zap.String("execution level", "Root"))

	userScriptAuthMiddleware := middleware.UserScriptAuthMiddleware(constants.USER_SCRIPT_AUTHENTICATION_ENDPOINT)
	go log.Info("authentication middleware initialization successful", zap.String("execution level", "Root"))

	userAuthMiddleware := middleware.UserAuthMiddleware(constants.USER_AUTHENTICATION_ENDPOINT)
	go log.Info("user authentication middleware initialization successful", zap.String("execution level", "Root"))

	router.Use(middleware.RequestLogger(log))
	go log.Info("registered logger for the router", zap.String("execution level", "Root"))

	router.Use(gin.Recovery())
	go log.Info("using panic recovery", zap.String("execution level", "Root"))

	routes.RegisterRoutes(router, log, db.DB, cache.RedisObj, authMiddleware, userScriptAuthMiddleware, userAuthMiddleware)

	go log.Info("Starting application at port 8080...", zap.String("Execution Level", "Root"))
	router.Run(":8080")
}
