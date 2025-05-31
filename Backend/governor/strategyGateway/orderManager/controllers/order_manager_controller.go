package controllers

import (
	"github.com/gin-gonic/gin"
	"governor/logger"
	"governor/strategyGateway/orderManager/services"
	"net/http"
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
