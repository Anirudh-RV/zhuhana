package main

import (
	"context"
	"governor/cache"
	"governor/db"
	"governor/kafka"
	"governor/kubernetescontroller"
	"governor/logger"
	"governor/middleware"
	"governor/routes"
	"governor/scheduler"

	constants "governor/constants"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	var ctx = context.Background()

	log := logger.NewLogger()
	go log.Info("logger started", zap.String("Execution Level", "Root"))

	db.InitDB(log)
	go log.Info("db connection successful", zap.String("Execution Level", "Root"))

	cache.InitRedis(ctx, log)
	go log.Info("redis connection successful", zap.String("Execution Level", "Root"))

	kubernetesService := kubernetescontroller.NewKubernetesService(log, db.DB)
	go log.Info("scheduler initialization successful", zap.String("Execution Level", "Root"))

	kafkaService := kafka.NewKafkaService(log, kubernetesService)
	kafkaService.Init(log)
	go log.Info("kafka initialization successful", zap.String("Execution Level", "Root"))

	schedulerService := scheduler.NewSchedulerService(cache.RedisObj, cache.RedisLockObj, log, db.DB, kafkaService)
	schedulerService.Init()
	go log.Info("scheduler initialization successful", zap.String("Execution Level", "Root"))

	router := gin.Default()
	go log.Info("Router setup successful", zap.String("Execution Level", "Root"))

	authMiddleware := middleware.AuthMiddleware(constants.API_AUTHENTICATION_ENDPOINT)
	go log.Info("authentication middleware initialization successful", zap.String("execution level", "Root"))

	userAuthMiddleware := middleware.UserAuthMiddleware(constants.USER_AUTHENTICATION_ENDPOINT)
	go log.Info("authentication middleware initialization successful", zap.String("execution level", "Root"))

	microserviceAuthenticator := middleware.NewMicroSeviceAuthenticator(log)
	microserviceAuthenticator.GetAllServiceTokens()
	go log.Info("microservice authenticator initialization successful", zap.String("execution level", "Root"))

	router.Use(middleware.RequestLogger(log))
	go log.Info("registered logger for the router", zap.String("execution level", "Root"))

	router.Use(gin.Recovery())
	go log.Info("using panic recovery", zap.String("execution level", "Root"))

	routes.RegisterRoutes(router, log, db.DB, cache.RedisObj, authMiddleware, userAuthMiddleware, microserviceAuthenticator, schedulerService, kafkaService)

	go log.Info("Starting application at port 8080...", zap.String("Execution Level", "Root"))
	router.Run(":8080")
}
