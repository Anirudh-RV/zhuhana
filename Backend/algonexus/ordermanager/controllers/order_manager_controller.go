package controllers

import (
	"algonexus/logger"
	"algonexus/ordermanager/models"
	"algonexus/ordermanager/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type OrderDomainHandlerFunc func(*models.OrderRequest) (*models.OrderResponse, error)

type OrderManagerController struct {
	logger   *logger.Logger
	service  *services.OrderManagerService
	handlers map[models.OrderDomain]OrderDomainHandlerFunc
}

func NewOrderManagerController(logger *logger.Logger, orderManagerService *services.OrderManagerService, handlers map[models.OrderDomain]OrderDomainHandlerFunc) *OrderManagerController {
	return &OrderManagerController{
		logger:   logger,
		service:  orderManagerService,
		handlers: handlers,
	}
}

func (omc *OrderManagerController) SubmitOrder(c *gin.Context) {

	var req models.Order

	if err := c.ShouldBindJSON(&req); err != nil {
		omc.logger.Error("failed to parse order request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	handler, ok := omc.handlers[req.Domain]

	if !ok {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "unsupported domain"})
		return
	}

	var orderRequest = &models.OrderRequest{
		Order:     req,
		OrderID:   uuid.New().String(),
		Timestamp: time.Now(),
	}

	res, err := handler(orderRequest)
	if err != nil {
		omc.logger.Error("handler error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}
