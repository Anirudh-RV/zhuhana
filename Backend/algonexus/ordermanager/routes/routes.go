package routes

import (
	logger "algonexus/logger"
	"algonexus/ordermanager/controllers"
	"algonexus/ordermanager/controllers/handlers"
	"algonexus/ordermanager/models"
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
	auth gin.HandlerFunc,
) {
	orderManagerRoutes := r.Group("ordermanager/")
	{
		orderManagerService := services.NewOrderManagerService(logger)
		go logger.Info("order manager service created", zap.String("execution level", "RegisterOrderManagerRoutesV1"))

		orderManagerControllerHandlers := map[models.OrderDomain]controllers.OrderDomainHandlerFunc{
			models.DomainBacktest: func(req *models.OrderRequest) (*models.OrderResponse, error) {
				return handlers.SubmitBacktestOrder(req)
			},
		}

		go logger.Info("order manager handler created", zap.String("execution level", "RegisterOrderManagerRoutesV1"))

		orderManagerController := controllers.NewOrderManagerController(
			logger, orderManagerService, orderManagerControllerHandlers)

		go logger.Info("order manager controller created", zap.String("execution level", "RegisterOrderManagerRoutesV1"))

		orderManagerRoutes.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})

		orderManagerRoutes.POST("/submit", orderManagerController.SubmitOrder)

	}

}
