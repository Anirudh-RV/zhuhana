package routes

import (
	"database/sql"
	"orchestrator/logger"

	_ "orchestrator/docs"

	outbound_handler_routes "orchestrator/outbound_handler/routes"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, log *logger.Logger, db *sql.DB, redis *redis.Client, authMiddleware gin.HandlerFunc) {
	v1 := r.Group("/api/outbound/v1/")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		outbound_handler_routes.OutboundHandlerRoutesV1(v1, log, db, redis)
	}
}
