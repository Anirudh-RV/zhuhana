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

func (oms *OrderManagerService) SubmitQueuedOrder(req *models.OrderRequest) (*models.OrderResponse, error) {
	//Round trip, sync (catch response) in service and return to controller

}
