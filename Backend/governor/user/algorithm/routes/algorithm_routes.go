package routes

import (
	"database/sql"
	"governor/logger"
	"governor/middleware"
	"governor/user/algorithm/controllers"
	"governor/user/algorithm/repositories"
	"governor/user/algorithm/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func UserAlgorithmRoutesV1(r *gin.RouterGroup, log *logger.Logger, db *sql.DB, redis *redis.Client, authMiddleware gin.HandlerFunc, userAuthMiddleware gin.HandlerFunc, microserviceAuthenticator *middleware.MicroSeviceAuthenticator) {
	algorithmRoutes := r.Group("user/algorithm/python/")
	{
		userAlgorithmRepository := repositories.NewUserAlgorithmRepository(db)
		go log.Info("user algorithm repository created", zap.String("execution level", "UserAlgorithmRoutesV1"))

		userAlgorithmService := services.NewUserAlgorithmService(log, userAlgorithmRepository, microserviceAuthenticator)
		go log.Info("upload script service created", zap.String("execution level", "UserAlgorithmRoutesV1"))

		userAlgorithmController := controllers.NewUserAlgorithmController(log, userAlgorithmService)
		go log.Info("upload script controller created", zap.String("execution level", "UserAlgorithmRoutesV1"))

		algorithmRoutes.POST("upload/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "header",
			Param:       "USER_TOKEN",
			EnableParam: true,
			Limit:       10,
			Window:      300,
			EnableIP:    true,
			IPLimit:     10,
			IPWindow:    300,
			Endpoint:    "/v1/user/algorith/script/upload/",
		}), userAuthMiddleware,
			userAlgorithmController.CreateUserAlgorithmHandler)
	}
}
