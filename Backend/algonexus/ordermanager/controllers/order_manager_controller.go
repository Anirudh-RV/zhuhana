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

	switch req.Domain {
	case models.DomainBacktest:
		omc.logger.Info("backtest order received")
		err = omc.SubmitBacktestOrder(orderRequest)
		omc.logger.Info("backtest order submitted")
	default:
		omc.logger.Info("invalid trade mode")
		c.JSON(http.StatusNotImplemented, gin.H{"error": "unsupported trade mode (has to be backtest/live/paper)"})
		return
	}

	if err != nil {
		omc.logger.Error("handler error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "success!"})
}

func (c *OrderManagerController) SubmitBacktestOrder(req *models.OrderRequest) error {
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

	err, channel := c.service.DeliverOrderToHub(req)

	if err != nil {
		c.logger.Error("order delivery to OrderHub failed", zap.Error(err))
		return err
	}
	c.logger.Info("order delivered to OrderHub")

	//go func(channel <-chan hubmodels.OrderEvent) {
	//	for event := range channel {
	//		//TODO Stop Condition
	//		s.logger.Info("OrderHub received event", zap.String("order_id", id), zap.String("type", string(event.Type)))
	//	}
	//}(channel)

	//!!! TEST ONLY !!!
	//go func() {
	//	event := <-channel
	//	c.logger.Info("OrderHub received event",
	//		zap.String("order_id", req.OrderID),
	//		zap.String("type", string(event.Type)),
	//	)
	//}()
	event := <-channel
	c.logger.Info("OrderHub received event",
		zap.String("order_id", req.OrderID),
		zap.String("type", string(event.Type)))
	
	return nil
}
