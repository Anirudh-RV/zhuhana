package services

import (
	logger "algonexus/logger"
	orderLogger "algonexus/ordermanager/logger"
	"algonexus/ordermanager/models"
)

type OrderManagerService struct {
	logger      *logger.Logger           //System-level Logger
	orderLogger *orderLogger.OrderLogger // Trade Information Logger
}

func NewOrderManagerService(logger *logger.Logger) *OrderManagerService {
	orderlogger := orderLogger.NewOrderLogger()
	return &OrderManagerService{
		logger:      logger,
		orderLogger: orderlogger,
	}
}

func (service *OrderManagerService) SubmitOrder(req models.OrderRequest) (*models.OrderResponse, error) {
	// TODO
	return nil, nil
}
