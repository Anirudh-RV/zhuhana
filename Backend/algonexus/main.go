package main

import (
	"algonexus/cache"
	"algonexus/constants"
	"algonexus/db"
	"algonexus/logger"
	"algonexus/middleware"
	"algonexus/ordermanager/backtestengine/broker"
	brokerservices "algonexus/ordermanager/backtestengine/broker/services"
	orderHubServices "algonexus/ordermanager/orderhub/services"
	"algonexus/routes"
	"context"
	"os"

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

	// Algonexus Service-level Infra init.
	// Broker first (in-process execution module), then inject it into OrderHub.
	executor := brokerservices.NewMockSimulator(log)
	brokerAdapter := broker.NewInProcessBroker(log, executor, broker.ParseOverflowPolicy(os.Getenv("BROKER_OVERFLOW_POLICY")))
	brokerAdapter.Start(ctx)
	go log.Info("Broker execution started", zap.String("Execution level", "Root"))

	orderHubService := orderHubServices.NewOrderHubService(log, brokerAdapter)
	go log.Info("OrderHub service started", zap.String("Execution level", "Root"))

	// HTTP logging is toggled off for load tests (GIN_MODE=release or HTTP_LOG=off):
	// gin.New() drops the per-request gin.Logger; we also skip the zap RequestLogger.
	var router *gin.Engine
	if os.Getenv("HTTP_LOG") == "off" || gin.Mode() == gin.ReleaseMode {
		router = gin.New()
		router.Use(gin.Recovery())
	} else {
		router = gin.Default() // gin.Logger + gin.Recovery
		router.Use(middleware.RequestLogger(log))
	}
	go log.Info("Router setup successful", zap.String("Execution Level", "Root"))

	microserviceAuthenticator := middleware.NewMicroSeviceAuthenticator(log)
	microserviceAuthenticator.GetAllServiceTokens()
	go log.Info("microservice authenticator initialization successful", zap.String("execution level", "Root"))

	authMiddleware := middleware.AuthMiddleware(constants.API_AUTHENTICATION_ENDPOINT)
	go log.Info("authentication middleware initialization successful", zap.String("execution level", "Root"))

	userAlgorithmAuthMiddleware := middleware.UserAlgorithmAuthMiddleware(constants.MICROSERVICE_USER_ALGORITHM_AUTHENTICATE_ENDPOINT, microserviceAuthenticator)
	go log.Info("user algorithm authentication middleware initialization successful", zap.String("execution level", "Root"))

	routes.RegisterRoutes(router, log, db.DB, &db.ClickHouse, cache.RedisObj, orderHubService, authMiddleware, userAlgorithmAuthMiddleware)

	go log.Info("Starting application at port 8080...", zap.String("Execution Level", "Root"))
	router.Run(":8080")
}
