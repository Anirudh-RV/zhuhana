package routes

import (
	"outbound/logger"

	_ "outbound/docs"

	stock_routes "outbound/marketdata/routes"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, log *logger.Logger, redis *redis.Client, authMiddleware gin.HandlerFunc) {
	v1 := r.Group("/api/marketdata/v1/")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		stock_routes.StocksRoutesV1(v1, log, redis)
	}
}
