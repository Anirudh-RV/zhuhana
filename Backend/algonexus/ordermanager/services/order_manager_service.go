package services

import (
	logger "algonexus/logger"
	orderLogger "algonexus/ordermanager/logger"
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
