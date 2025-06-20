package controllers

import (
	"algonexus/logger"
	"algonexus/ordermanager/controllers/orderfsm"
	"algonexus/ordermanager/models"
	"algonexus/ordermanager/services"
	"context"
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
		Order:           req,
		OrderID:         uuid.New().String(),
		Timestamp:       time.Now(),
		ResponseChannel: make(chan models.OrderResponse, 1),
	}

	var err error = nil

	orderFSM := orderfsm.NewOrderFSM(orderRequest)

	switch req.Domain {
	case models.DomainBacktest:
		omc.logger.Info("backtest order received")
		err = omc.HandleBacktestOrder(orderFSM)
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

func (c *OrderManagerController) HandleBacktestOrder(order *orderfsm.OrderFSM) error {
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

	err := order.Transition(models.StatusPendingSend)
	if err != nil {
		return err
	}

	c.logger.Info("order pending to submit to queue")

	err = c.service.EnqueueOrder(order.OrderRequest)

	if err != nil {
		c.logger.Error("order enqueue failed", zap.Error(err))
		return err
	}

	err = order.Transition(models.StatusEnqueued)
	if err != nil {
		return err
	}

	c.logger.Info("order submitted to queue")

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		select {
		case resp := <-order.OrderRequest.ResponseChannel:
			c.logger.Info("order was consumed",
				zap.String("order ID", resp.OrderID),
				zap.String("status", string(resp.Status)))

			err := order.Transition(resp.Status)
			if err != nil {
				c.logger.Error("Order Status Transition Error")
				return
			}
		case <-ctx.Done(): //timeout
			return
		default:
			c.logger.Warning("Response channel unavailable or closed (timeout)")
		}
	}()

	return nil
}
