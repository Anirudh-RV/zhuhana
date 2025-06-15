package controllers

import (
	"algonexus/logger"
	"algonexus/orderManager/services"
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
