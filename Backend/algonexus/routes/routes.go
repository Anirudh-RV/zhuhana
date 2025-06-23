package routes

import (
	_ "algonexus/docs"
	"algonexus/logger"
	backtestRoutes "algonexus/ordermanager/backtestengine/routes"
	orderHubServices "algonexus/ordermanager/orderhub/services"
	orderManagerRoutes "algonexus/ordermanager/routes"
	"database/sql"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, log *logger.Logger, db *sql.DB, clickHouse *clickhouse.Conn, redis *redis.Client, orderHubService *orderHubServices.OrderHubService, authMiddleware gin.HandlerFunc, userAlgorithmAuthMiddleware gin.HandlerFunc) {

	v1 := r.Group("/v1/")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})
		orderManagerRoutes.RegisterOrderManagerRoutesV1(v1, log, db, clickHouse, redis, orderHubService, authMiddleware)
		backtestRoutes.RegisterBacktestRoutesV1(v1, log, db, clickHouse, redis, authMiddleware, userAlgorithmAuthMiddleware)
	}
}
