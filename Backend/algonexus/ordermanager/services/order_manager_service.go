package services

import (
	EQservices "algonexus/eventqueue/services"
	logger "algonexus/logger"
	orderLogger "algonexus/ordermanager/logger"
	"algonexus/ordermanager/models"
	"context"
)

type OrderManagerService struct {
	logger         *logger.Logger           //System-level Logger
	orderLogger    *orderLogger.OrderLogger // Trade Information Logger
	rsOrderService *EQservices.RsOrderService
}

func NewOrderManagerService(logger *logger.Logger, orderservice *EQservices.RsOrderService) *OrderManagerService {
	orderlogger := orderLogger.NewOrderLogger()
	return &OrderManagerService{
		logger:         logger,
		orderLogger:    orderlogger,
		rsOrderService: orderservice,
	}
}

func (oms *OrderManagerService) EnqueueOrder(req *models.OrderRequest) error {
	var ctx = context.Background()
	//Round trip, sync (catch response) in service and return to controller

	err := oms.rsOrderService.PushOrder(ctx, *req)
	if err != nil {
		return err
	}
	return nil
}
