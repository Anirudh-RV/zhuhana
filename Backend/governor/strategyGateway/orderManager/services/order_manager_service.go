package services

import (
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	governorLogger "governor/logger"
	orderManagerLogger "governor/strategyGateway/orderManager/logger"
	"governor/strategyGateway/orderManager/models"
	"governor/strategyGateway/orderManager/services/handler"
	"time"
)

type OrderManagerService struct {
	logger         *governorLogger.Logger     //System-level Logger
	orderLogger    *orderManagerLogger.Logger // Trade Information Logger
	domainHandlers map[models.OrderDomain]func(*models.OrderRequest) (*models.OrderResponse, error)
}

func NewOrderManagerService(logger *governorLogger.Logger, orderLogger *orderManagerLogger.Logger) *OrderManagerService {
	return &OrderManagerService{
		logger:      logger,
		orderLogger: orderLogger,
		domainHandlers: map[models.OrderDomain]func(*models.OrderRequest) (*models.OrderResponse, error){
			models.DomainBacktest: handler.HandleBacktestOrder,
			// More Domains/Exchanges/Brokers
		},
	}
}

func (s *OrderManagerService) SubmitOrder(req *models.OrderRequest) (*models.OrderResponse, error) {
	//if err := req.Validate(); err != nil {
	//	return nil, fmt.Errorf("validation failed: %w", err)
	//}

	// ======= Backtest For Now =======
	var order = req.Order

	orderHandler, ok := s.domainHandlers[order.Domain]
	if !ok {
		s.logger.Error("unsupported order domain")
		return nil, fmt.Errorf("unsupported order domain")
	}

	var orderRequest = &models.OrderRequest{
		Order:     order,
		OrderID:   uuid.New().String(),
		Timestamp: time.Now(),
	}

	response, err := orderHandler(orderRequest)

	if err != nil {
		s.logger.Error("Error occurs in order request submission", zap.Error(err))
		return nil, err
	}

	return response, nil

}
