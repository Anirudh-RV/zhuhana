package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"governor/logger"
	"governor/strategyGateway/orderManager/models"
	"governor/strategyGateway/orderManager/services"
	"net/http"
	"time"
)

type OrderManagerController struct {
	log                 *logger.Logger
	orderManagerService *services.OrderManagerService
}

func NewOrderManagerController(log *logger.Logger, orderManagerService *services.OrderManagerService) *OrderManagerController {
	return &OrderManagerController{
		log:                 log,
		orderManagerService: orderManagerService,
	}
}

func (omc *OrderManagerController) HandleSubmitOrder(c *gin.Context) {
	var order models.Order
	if err := json.NewDecoder(c.Request.Body).Decode(&order); err != nil {
		c.JSON(http.StatusBadRequest, models.OrderResponse{
			Status:  string(models.ResponseStatusError),
			Message: "Invalid request payload",
		})
		return
	}

	var orderRequest models.OrderRequest

	orderRequest.Order = order
	orderRequest.Timestamp = time.Now()
	orderRequest.OrderID = uuid.New().String()

	// WebSocket Streaming
	resp, err := omc.orderManagerService.SubmitOrder(&orderRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.OrderResponse{
			Status:  string(models.ResponseStatusError),
			Message: "Internal Server Error",
		})
		return
	}

	resp.Status = string(models.ResponseStatusSubmitted)
	resp.Message = "Order successfully submitted"
	c.JSON(http.StatusCreated, resp)

}
