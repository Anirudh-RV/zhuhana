package routes

import (
	EQservices "algonexus/eventqueue/services"
	logger "algonexus/logger"
	"algonexus/ordermanager/controllers"
	"algonexus/ordermanager/services"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"

	"database/sql"

	"github.com/redis/go-redis/v9"
)

func RegisterOrderManagerRoutesV1(
	r *gin.RouterGroup,
	logger *logger.Logger,
	db *sql.DB,
	redis *redis.Client,
	rsOrderService *EQservices.RsOrderService,
	auth gin.HandlerFunc,
) {

	orderManagerRoutes := r.Group("ordermanager/")
	{
		orderManagerService := services.NewOrderManagerService(logger, rsOrderService)
		go logger.Info("order manager service created", zap.String("execution level", "RegisterOrderManagerRoutesV1"))

		orderManagerController := controllers.NewOrderManagerController(
			logger, orderManagerService)

		go logger.Info("order manager controller created", zap.String("execution level", "RegisterOrderManagerRoutesV1"))

		orderManagerRoutes.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})

		orderManagerRoutes.POST("/submit", orderManagerController.SubmitOrder)

	}

}
