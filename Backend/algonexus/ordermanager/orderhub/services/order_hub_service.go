package services

import (
	eqServices "algonexus/eventqueue/services"
	"algonexus/logger"
	"algonexus/ordermanager/models"
	"algonexus/ordermanager/orderhub/runtime"
	"context"
	"go.uber.org/zap"
	"sync"
)

type OrderHubService struct {
	logger         *logger.Logger
	orders         map[string]*runtime.OrderHandle
	rsOrderService *eqServices.RsOrderService
	rwmu           sync.RWMutex
}

func NewOrderHubService(logger *logger.Logger) *OrderHubService {
	return &OrderHubService{
		logger: logger,
		orders: make(map[string]*runtime.OrderHandle),
	}
}

func (s *OrderHubService) RegisterOrder(req *models.OrderRequest) {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()
	handle := runtime.NewOrderHandle(req)
	s.orders[req.OrderID] = handle

	// Producer
	go func() {
		err := s.rsOrderService.PushOrder(context.Background(), req)
		if err != nil {
			s.logger.Error("unable to push order request into the event queue", zap.Error(err))
			return
		}

		s.logger.Info("order successfully enqueued", zap.String("request", req.OrderID))

		err = handle.OrderFlow.Transition(models.StatusEnqueued)
		if err != nil {
			s.logger.Error("transition failed", zap.String("request", req.OrderID), zap.Error(err))
			return
		}
	}()

	// Listener routine (listen to consumer)
	go s.Listen(req.OrderID)

}

func (s *OrderHubService) Listen(id string) {
	s.rwmu.RLock()
	handle, ok := s.orders[id]
	s.rwmu.RUnlock()

	if !ok {
		s.logger.Warning("attempted to listen to unknown order", zap.String("id", id))
		return
	}

	for event := range handle.Channel {
		s.logger.Info("received event", zap.String("order_id", id), zap.String("type", string(event.Type)))
	}
}

func (s *OrderHubService) UnregisterOrder(id string) {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	orderHandle, ok := s.orders[id]

	if !ok {
		s.logger.Warning("can't find order", zap.String("orderID", id))
		return
	}

	if !orderHandle.OrderFlow.IsTerminated() {
		s.logger.Fatal("Order delete too early!", zap.String("orderID", id))
		return
	}

	delete(s.orders, id)

}
