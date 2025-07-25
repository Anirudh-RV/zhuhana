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

func UserSecretsRoutesV1(r *gin.RouterGroup, log *logger.Logger, db *sql.DB, redis *redis.Client, authMiddleware gin.HandlerFunc, userScriptAuthMiddleware gin.HandlerFunc, userAuthMiddleware gin.HandlerFunc) {
	userSecretsRepoObj := repositories.NewUserSecretRepository(db)
	go log.Info("user secrets service created", zap.String("execution level", "UserSecretsRoutesV1"))

	userSecretsServiceObj := services.NewUserSecretsService(log, userSecretsRepoObj)
	go log.Info("user secrets service created", zap.String("execution level", "UserSecretsRoutesV1"))

	userSecretsSetController := controllers.NewUserSecretsSetController(userSecretsServiceObj, log)
	go log.Info("user secrets set controller created", zap.String("execution level", "UserSecretsRoutesV1"))

	userSecretsGetController := controllers.NewUserSecretsGetController(userSecretsServiceObj, log)
	go log.Info("user secrets get controller created", zap.String("execution level", "UserSecretsRoutesV1"))

	userSecretsDeleteController := controllers.NewUserSecretsDeleteController(userSecretsServiceObj, log)
	go log.Info("user secrets delete controller created", zap.String("execution level", "UserSecretsRoutesV1"))

	user := r.Group("user/")
	{
		user.POST("secret/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "header",
			Param:       "USER_TOKEN",
			EnableParam: true,
			Limit:       300,
			Window:      300,
			EnableIP:    false,
			Endpoint:    "POST/user/secret/",
		}), userAuthMiddleware,
			userSecretsSetController.UserSecretsSetHandler)

		user.GET("secret/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "header",
			Param:       "USER_TOKEN",
			EnableParam: true,
			Limit:       300,
			Window:      300,
			EnableIP:    false,
			Endpoint:    "GET/user/secret/",
		}), userAuthMiddleware,
			userSecretsGetController.UserSecretGetHandler)

		user.DELETE("secret/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "header",
			Param:       "USER_TOKEN",
			EnableParam: true,
			Limit:       300,
			Window:      300,
			EnableIP:    false,
			Endpoint:    "DELETE/user/secret/",
		}), userAuthMiddleware,
			userSecretsDeleteController.UserSecretDeleteHandler)

		user.GET("secret/keys/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "header",
			Param:       "USER_TOKEN",
			EnableParam: true,
			Limit:       300,
			Window:      300,
			EnableIP:    false,
			Endpoint:    "/user/secret/keys",
		}), userAuthMiddleware,
			userSecretsGetController.UserSecretKeysGetHandler)
	}

}
