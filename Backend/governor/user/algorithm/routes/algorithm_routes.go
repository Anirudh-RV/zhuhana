package routes

import (
	"database/sql"
	"governor/kafka"
	"governor/kubernetescontroller"
	"governor/logger"
	"governor/middleware"
	"governor/scheduler"
	"governor/user/algorithm/controllers"
	"governor/user/algorithm/repositories"
	"governor/user/algorithm/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func UserAlgorithmRoutesV1(r *gin.RouterGroup, log *logger.Logger, db *sql.DB, redis *redis.Client, authMiddleware gin.HandlerFunc, userAuthMiddleware gin.HandlerFunc, userAlgorithmAuthMiddleware gin.HandlerFunc, microserviceAuthenticator *middleware.MicroSeviceAuthenticator, schedulerService *scheduler.SchedulerService, kafkaService *kafka.KafkaService, kubernetesService *kubernetescontroller.KubernetesService) {
	algorithmRoutes := r.Group("user/algorithm/")
	{
		userAlgorithmRepository := repositories.NewUserAlgorithmRepository(db)
		go log.Info("user algorithm repository created", zap.String("execution level", "UserAlgorithmRoutesV1"))

		userAlgorithmService := services.NewUserAlgorithmService(log, userAlgorithmRepository, microserviceAuthenticator, schedulerService, kafkaService, kubernetesService)
		go log.Info("upload script service created", zap.String("execution level", "UserAlgorithmRoutesV1"))

		userAlgorithmController := controllers.NewUserAlgorithmController(log, userAlgorithmService)
		go log.Info("upload script controller created", zap.String("execution level", "UserAlgorithmRoutesV1"))

		algorithmRoutes.POST("schedule/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "body",
			Param:       "algorithmID",
			EnableParam: true,
			Limit:       10,
			Window:      300,
			EnableIP:    true,
			IPLimit:     10,
			IPWindow:    300,
			Endpoint:    "/v1/user/algorithm/schedule/",
		}), userAuthMiddleware,
			userAlgorithmController.UpdateUserAlgorithmCronSchedule)

		algorithmRoutes.POST("start/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "body",
			Param:       "algorithmID",
			EnableParam: true,
			Limit:       10,
			Window:      300,
			EnableIP:    true,
			IPLimit:     10,
			IPWindow:    300,
			Endpoint:    "/v1/user/algorithm/start/",
		}), userAuthMiddleware,
			userAlgorithmController.StartUserAlgorithm)

		algorithmRoutes.POST("stop/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "body",
			Param:       "algorithmID",
			EnableParam: true,
			Limit:       10,
			Window:      300,
			EnableIP:    true,
			IPLimit:     10,
			IPWindow:    300,
			Endpoint:    "/v1/user/algorithm/stop/",
		}), userAuthMiddleware,
			userAlgorithmController.StopUserAlgorithm)

		algorithmRoutes.POST("schedule/cancel/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "body",
			Param:       "algorithmID",
			EnableParam: true,
			Limit:       10,
			Window:      300,
			EnableIP:    true,
			IPLimit:     10,
			IPWindow:    300,
			Endpoint:    "/v1/user/algorithm/schedule/cancel/",
		}), userAuthMiddleware,
			userAlgorithmController.CancelUserAlgorithmCronSchedule)

		algorithmRoutes.GET("info/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "header",
			Param:       "USER_TOKEN",
			EnableParam: true,
			Limit:       100,
			Window:      300,
			EnableIP:    true,
			IPLimit:     100,
			IPWindow:    300,
			Endpoint:    "/v1/user/algorithm/info/",
		}), userAuthMiddleware,
			userAlgorithmController.GetUserAlgorithmByID)

		algorithmRoutes.GET("/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "header",
			Param:       "USER_TOKEN",
			EnableParam: true,
			Limit:       300,
			Window:      300,
			EnableIP:    true,
			IPLimit:     300,
			IPWindow:    300,
			Endpoint:    "/v1/user/algorithm/",
		}), userAuthMiddleware,
			userAlgorithmController.GetUserAlgorithms)

		pythonAlgorithms := algorithmRoutes.Group("python/")
		{
			pythonAlgorithms.POST("upload/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
				Source:      "header",
				Param:       "USER_TOKEN",
				EnableParam: true,
				Limit:       10,
				Window:      300,
				EnableIP:    true,
				IPLimit:     10,
				IPWindow:    300,
				Endpoint:    "/v1/user/algorith/python/upload/",
			}), userAuthMiddleware,
				userAlgorithmController.CreateUserAlgorithmHandler)

			pythonAlgorithms.PUT("edit/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
				Source:      "header",
				Param:       "USER_TOKEN",
				EnableParam: true,
				Limit:       10,
				Window:      300,
				EnableIP:    true,
				IPLimit:     10,
				IPWindow:    300,
				Endpoint:    "/v1/user/algorith/python/edit/",
			}), userAuthMiddleware,
				userAlgorithmController.EditUserAlgorithmHandler)
		}
	}
}
