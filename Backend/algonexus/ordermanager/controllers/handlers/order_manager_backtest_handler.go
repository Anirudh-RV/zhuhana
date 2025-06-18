package handlers

import (
	"algonexus/ordermanager/models"
	"github.com/google/uuid"
	"time"
)

func SubmitBacktestOrder(req *models.OrderRequest) (*models.OrderResponse, error) {
	return &models.OrderResponse{
		OrderID:       req.OrderID,
		OrderDetails:  req.Order,
		SubmitTime:    req.Timestamp,
		BrokerOrderID: "SIM-" + uuid.New().String(),
		Status:        models.ResponseStatusSubmitted,
		Message:       "Order successfully accepted in simulation.",
		Fills:         []models.OrderFill{},
		Time:          time.Now(),
	}, nil
}
