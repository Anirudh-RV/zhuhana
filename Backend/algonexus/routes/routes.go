package routes

import (
	_ "algonexus/docs"
	"algonexus/eventqueue"
	logger "algonexus/logger"
	"database/sql"

	ordermanagerroutes "algonexus/ordermanager/routes"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, log *logger.Logger, db *sql.DB, redis *redis.Client, rsOrderService *eventqueue.RsOrderService, authMiddleware gin.HandlerFunc) {

	v1 := r.Group("/v1/algonexus")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})
		ordermanagerroutes.RegisterOrderManagerRoutesV1(v1, log, db, redis, rsOrderService, authMiddleware)
	}
}
