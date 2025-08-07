package main

import (
	"algonexus/cache"
	"algonexus/constants"
	"algonexus/db"
	"algonexus/logger"
	"algonexus/middleware"
	orderHubServices "algonexus/ordermanager/orderhub/services"
	"algonexus/routes"
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	var ctx = context.Background()

	log := logger.NewLogger()
	go log.Info("Logger started", zap.String("Execution Level", "Root"))

	db.InitDB(log)
	go log.Info("Postgres DB connection successful", zap.String("Execution Level", "Root"))

	db.InitClickHouse(log)
	go log.Info("ClickHouse DB connection successful", zap.String("Execution Level", "Root"))

	cache.InitRedis(ctx, log)
	go log.Info("Redis connection successful", zap.String("Execution Level", "Root"))

	// Algonexus Service-level Infra init
	orderHubService := orderHubServices.NewOrderHubService(log)
	go log.Info("OrderHub service started", zap.String("Execution level", "Root"))

	//router := gin.Default()
	router := gin.New()
	go log.Info("Router setup successful", zap.String("Execution Level", "Root"))

	microserviceAuthenticator := middleware.NewMicroSeviceAuthenticator(log)
	microserviceAuthenticator.GetAllServiceTokens()
	go log.Info("microservice authenticator initialization successful", zap.String("execution level", "Root"))

	authMiddleware := middleware.AuthMiddleware(constants.API_AUTHENTICATION_ENDPOINT)
	go log.Info("authentication middleware initialization successful", zap.String("execution level", "Root"))

	userAlgorithmAuthMiddleware := middleware.UserAlgorithmAuthMiddleware(constants.MICROSERVICE_USER_ALGORITHM_AUTHENTICATE_ENDPOINT, microserviceAuthenticator)
	go log.Info("user algorithm authentication middleware initialization successful", zap.String("execution level", "Root"))

	router.Use(middleware.RequestLogger(log))
	go log.Info("registered logger for the router", zap.String("execution level", "Root"))

	router.Use(gin.Recovery())
	go log.Info("using panic recovery", zap.String("execution level", "Root"))

	routes.RegisterRoutes(router, log, db.DB, &db.ClickHouse, cache.RedisObj, orderHubService, authMiddleware, userAlgorithmAuthMiddleware)

	go log.Info("Starting application at port 8080...", zap.String("Execution Level", "Root"))
	router.Run(":8080")
}
