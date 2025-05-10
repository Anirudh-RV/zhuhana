package routes

import (
	"context"
	"database/sql"
	"uasam/commonutils"
	"uasam/email"
	"uasam/logger"

	_ "uasam/docs"

	microservice_routes "uasam/microservices/routes"
	user_routes "uasam/users/routes"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(ctx *context.Context, r *gin.Engine, log *logger.Logger, db *sql.DB, redis *redis.Client, emailService *email.EmailService, jwtService *commonutils.JWTService) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := r.Group("/v1/")
	{
		user_routes.UserRoutesV1(ctx, v1, log, db, redis, emailService, jwtService)
		microservice_routes.MicroServiceRoutesV1(v1, log, db, redis, jwtService)
	}
}
