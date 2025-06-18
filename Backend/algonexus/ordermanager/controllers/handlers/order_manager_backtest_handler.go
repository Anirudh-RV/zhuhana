package handlers

import (
	"algonexus/ordermanager/models"
)

func SubmitBacktestOrder(*models.OrderRequest) (*models.OrderResponse, error) {
	return &models.OrderResponse{}, nil
}
