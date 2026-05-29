package services

import (
	"algonexus/logger"
	orderLogger "algonexus/ordermanager/logger"
	"algonexus/ordermanager/models"
	orderHubServices "algonexus/ordermanager/orderhub/services"
)

type OrderManagerService struct {
	logger      *logger.Logger           //System-level Logger
	orderLogger *orderLogger.OrderLogger // Trade Information Logger

	orderHubService *orderHubServices.OrderHubService
}

func NewOrderManagerService(logger *logger.Logger, hubService *orderHubServices.OrderHubService) *OrderManagerService {
	orderlogger := orderLogger.NewOrderLogger()
	return &OrderManagerService{
		logger:          logger,
		orderLogger:     orderlogger,
		orderHubService: hubService,
	}
}

func (oms *OrderManagerService) DeliverOrderToHub(req *models.OrderRequest) error {
	//Register Order Handle in hub first
	oms.orderHubService.RegisterOrder(req)
	return nil
}
