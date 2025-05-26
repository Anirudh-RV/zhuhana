package routes

import (
	"database/sql"
	"secretsmanager/logger"

	_ "secretsmanager/docs"

	usersecrets_routes "secretsmanager/secrets/usersecrets/routes"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, log *logger.Logger, db *sql.DB, redis *redis.Client, authMiddleware gin.HandlerFunc, userAuthMiddleware gin.HandlerFunc) {
	v1 := r.Group("/v1/")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		usersecrets_routes.UserSecretsRoutesV1(v1, log, db, redis, authMiddleware, userAuthMiddleware)
	}
}
