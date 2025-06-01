package services

import (
	governorLogger "governor/logger"
	orderManagerLogger "governor/strategyGateway/orderManager/logger"
	"governor/strategyGateway/orderManager/models"
)

type OrderManagerService struct {
	logger      *governorLogger.Logger     //System-level Logger
	orderLogger *orderManagerLogger.Logger // Trade Information Logger
}

func NewOrderManagerService(logger *governorLogger.Logger, orderLogger *orderManagerLogger.Logger) *OrderManagerService {
	return &OrderManagerService{
		logger:      logger,
		orderLogger: orderLogger,
	}
}

func (service *OrderManagerService) SubmitOrder(req models.OrderRequest) (*models.OrderResponse, error) {
	// TODO
	return nil, nil
}
