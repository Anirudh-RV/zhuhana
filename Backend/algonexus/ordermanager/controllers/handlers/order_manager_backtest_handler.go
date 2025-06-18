package handlers

import (
	"algonexus/ordermanager/models"
	"github.com/google/uuid"
	"time"
)

func SubmitBacktestOrder(req *models.OrderRequest) (*models.OrderResponse, error) {
	return &models.OrderResponse{
		OrderID:       req.OrderID,
		BrokerOrderID: "SIM-" + uuid.New().String(),
		Status:        models.ResponseStatusSubmitted,
		Message:       "Order successfully accepted in simulation.",
		Fills:         []models.OrderFill{},
		Time:          time.Now(),
	}, nil
}
