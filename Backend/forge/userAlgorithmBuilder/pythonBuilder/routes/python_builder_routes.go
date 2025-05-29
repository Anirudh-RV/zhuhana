package routes

import (
	"database/sql"
	"forge/dockercontroller"
	"forge/logger"
	"forge/middleware"
	"forge/userAlgorithmBuilder/pythonBuilder/controllers"
	"forge/userAlgorithmBuilder/pythonBuilder/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func PythonBuilderRoutesV1(r *gin.RouterGroup, log *logger.Logger, db *sql.DB, redis *redis.Client, authMiddleware gin.HandlerFunc, dockerService *dockercontroller.DockerService) {
	pythonBuilderRoutes := r.Group("algorithm/python/")
	{
		pythonBuilderService := services.NewPythonBuilderService(log, dockerService)
		go log.Info("python builder service created", zap.String("execution level", "PythonBuilderRoutesV1"))

		pythonBuilderController := controllers.NewPythonBuilderController(log, pythonBuilderService)
		go log.Info("python builder controller created", zap.String("execution level", "PythonBuilderRoutesV1"))

		pythonBuilderRoutes.POST("build/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "body",
			Param:       "userID",
			EnableParam: true,
			Limit:       300,
			Window:      300,
			EnableIP:    false,
			Endpoint:    "/algorithm/python/build/",
		}), authMiddleware,
			pythonBuilderController.PythonBuilderHandler)
	}

}
