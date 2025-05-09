package routes

import (
	"database/sql"
	"uasam/logger"

	_ "uasam/docs"

	microservice_routes "uasam/microservices/routes"
	user_routes "uasam/users/routes"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, log *logger.Logger, db *sql.DB, redis *redis.Client) {
	v1 := r.Group("/api/users/v1/")
	{
		v1.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		user_routes.UserRoutesV1(v1, log, db, redis)
		microservice_routes.MicroServiceRoutesV1(v1, log, db, redis)
	}
}
