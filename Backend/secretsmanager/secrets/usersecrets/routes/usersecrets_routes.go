package routes

import (
	"database/sql"
	"secretsmanager/logger"
	"secretsmanager/middleware"
	"secretsmanager/secrets/usersecrets/controllers"
	"secretsmanager/secrets/usersecrets/repositories"
	"secretsmanager/secrets/usersecrets/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func UserSecretsRoutesV1(r *gin.RouterGroup, log *logger.Logger, db *sql.DB, redis *redis.Client, authMiddleware gin.HandlerFunc, userAuthMiddleware gin.HandlerFunc) {
	userSecretsRepoObj := repositories.NewUserSecretRepository(db)
	go log.Info("user secrets service created", zap.String("execution level", "UserSecretsRoutesV1"))

	userSecretsServiceObj := services.NewUserSecretsService(log, userSecretsRepoObj)
	go log.Info("user secrets service created", zap.String("execution level", "UserSecretsRoutesV1"))

	userSecretsSetController := controllers.NewUserSecretsSetController(userSecretsServiceObj, log)
	go log.Info("user secrets set controller created", zap.String("execution level", "UserSecretsRoutesV1"))

	r.POST("user/secrets/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
		Source:      "header",
		Param:       "USER_TOKEN",
		EnableParam: true,
		Limit:       300,
		Window:      300,
		EnableIP:    false,
		Endpoint:    "POST/user/secrets/",
	}), authMiddleware,
		userAuthMiddleware,
		userSecretsSetController.UserSecretsSetHandler)

	r.GET("user/secrets/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
		Source:      "header",
		Param:       "USER_TOKEN",
		EnableParam: true,
		Limit:       300,
		Window:      300,
		EnableIP:    false,
		Endpoint:    "GET/user/secrets/",
	}), authMiddleware,
		userAuthMiddleware,
		userSecretsSetController.UserSecretsGetHandler)

}
