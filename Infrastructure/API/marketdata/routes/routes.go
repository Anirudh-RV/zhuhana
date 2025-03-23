package routes

import (
	"marketdata/logger"

	_ "marketdata/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, log *logger.Logger) {
	v1 := r.Group("/api/marketdata/v1/")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		StocksRoutesV1(v1, log)
	}
}
