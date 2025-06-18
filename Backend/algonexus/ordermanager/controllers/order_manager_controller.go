package controllers

import (
	"algonexus/logger"
	"algonexus/ordermanager/models"
	"algonexus/ordermanager/services"
)

type OrderDomainHandlerFunc func(*models.OrderRequest) (*models.OrderResponse, error)

type OrderManagerController struct {
	logger   *logger.Logger
	service  *services.OrderManagerService
	handlers map[models.OrderDomain]OrderDomainHandlerFunc
}

func NewOrderManagerController(logger *logger.Logger, orderManagerService *services.OrderManagerService, handlers map[models.OrderDomain]OrderDomainHandlerFunc) *OrderManagerController {
	return &OrderManagerController{
		logger:   logger,
		service:  orderManagerService,
		handlers: handlers,
	}
}

func (omc *OrderManagerController) SubmitOrder(req *models.OrderRequest) (*models.OrderResponse, error) {
	//if err := req.Validate(); err != nil {
	//	return nil, fmt.Errorf("validation failed: %w", err)
	//}

	// ======= Backtest For Now =======
	//var order = req.Order
	//
	//var orderRequest = &models.OrderRequest{
	//	Order:     order,
	//	OrderID:   uuid.New().String(),
	//	Timestamp: time.Now(),
	//}

	//if err != nil {
	//	s.logger.Error("Error occurs in order request submission", zap.Error(err))
	//	return nil, err
	//}

	return nil, nil
}
