package routes

import (
	"database/sql"
	"governor/kafka"
	"governor/logger"
	"governor/middleware"
	"governor/scheduler"

	_ "governor/docs"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	StradegyGatewayRoutesV1 "governor/strategyGateway/routes"
	UserUserAlgorithmRoutesV1 "governor/user/algorithm/routes"
)

func RegisterRoutes(r *gin.Engine, log *logger.Logger, db *sql.DB, redis *redis.Client, authMiddleware gin.HandlerFunc, userAuthMiddleware gin.HandlerFunc, microserviceAuthenticator *middleware.MicroSeviceAuthenticator, schedulerService *scheduler.SchedulerService, kafkaService *kafka.KafkaService) {
	v1 := r.Group("/v1/")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})

		// Register other routes here
		StradegyGatewayRoutesV1.RegisterStrategyGatewayRoutesV1(v1, log, db, redis, authMiddleware)
		UserUserAlgorithmRoutesV1.UserAlgorithmRoutesV1(v1, log, db, redis, authMiddleware, userAuthMiddleware, microserviceAuthenticator, schedulerService, kafkaService)

	}
}
