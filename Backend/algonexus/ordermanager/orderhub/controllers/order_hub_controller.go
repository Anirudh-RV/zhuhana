package controllers

import (
	"algonexus/logger"
	"algonexus/ordermanager/orderhub/services"
)

type OrderHubController struct {
	logger  *logger.Logger
	service *services.OrderHubService
}

func NewOrderHubController(logger *logger.Logger, service *services.OrderHubService) *OrderHubController {
	return &OrderHubController{
		logger:  logger,
		service: service,
	}
}

//func (c *OrderHubController) AddOrder(req *models.OrderRequest) {
//	//TODO More request validation
//	if req == nil {
//		c.logger.Warning("received nil order request")
//		return
//	}
//
//	c.service.RegisterOrder(req)
//}
