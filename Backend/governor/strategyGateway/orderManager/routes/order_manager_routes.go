package routes

import (
	"github.com/gin-gonic/gin"
	"governor/logger"

	"database/sql"
	"github.com/redis/go-redis/v9"
)

func RegisterTradeOrderManagerRoutesV1(
	r *gin.RouterGroup,
	log *logger.Logger,
	db *sql.DB,
	redis *redis.Client,
	auth gin.HandlerFunc,
) {
	orderManager := r.Group("/order")
	// orderManager.Use(auth)
	// orderManager.POST("", handler.SubmitOrder(db))
	orderManager.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
}
