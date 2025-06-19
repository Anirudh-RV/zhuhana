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

type OrderHandlerFunc func(*models.OrderRequest) (*models.OrderResponse, error)

type OrderManagerController struct {
	logger  *logger.Logger
	service *services.OrderManagerService
}

func NewOrderManagerController(logger *logger.Logger, orderManagerService *services.OrderManagerService) *OrderManagerController {
	return &OrderManagerController{
		logger:  logger,
		service: orderManagerService,
	}
}

func (omc *OrderManagerController) SubmitOrder(c *gin.Context) {

	var req models.Order

	if err := c.ShouldBindJSON(&req); err != nil {
		omc.logger.Error("failed to parse order request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	var orderRequest = &models.OrderRequest{
		Order:     req,
		OrderID:   uuid.New().String(),
		Timestamp: time.Now(),
	}

	var err error = nil
	var response *models.OrderResponse = nil

	switch req.Domain {
	case models.DomainBacktest:
		omc.logger.Info("backtest order received")
		response, err = omc.SubmitBacktestOrder(orderRequest)
		omc.logger.Info("backtest order submitted")
	default:
		omc.logger.Info("invalid trade mode")
		c.JSON(http.StatusNotImplemented, gin.H{"error": "unsupported trade mode (has to be backtest/live/paper)"})
		return
	}

	if err != nil || response == nil {
		omc.logger.Error("handler error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (c *OrderManagerController) SubmitBacktestOrder(req *models.OrderRequest) (*models.OrderResponse, error) {
	//return &models.OrderResponse{
	//	OrderID:       req.OrderID,
	//	OrderDetails:  req.Order,
	//	SubmitTime:    req.Timestamp,
	//	BrokerOrderID: "SIM-" + uuid.New().String(),
	//	Status:        models.ResponseStatusSubmitted,
	//	Message:       "Order successfully accepted in simulation.",
	//	Fills:         []models.OrderFill{},
	//	Time:          time.Now(),
	//}, nil

	c.logger.Info("order pending to submit to queue")
	response, err := c.service.SubmitQueuedOrder(req)

	if err != nil {
		c.logger.Error("order submitted failed", zap.Error(err))
		return nil, err
	}

	c.logger.Info("order submitted to queue")
	
	return response, nil
}
