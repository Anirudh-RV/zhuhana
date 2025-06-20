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
	Logger         *logger.Logger
	Orders         map[string]*runtime.OrderHandle
	rsOrderService *eqServices.RsOrderService
	RWmu           sync.RWMutex
}

func NewOrderHubService(logger *logger.Logger) *OrderHubService {
	return &OrderHubService{
		Logger: logger,
		Orders: make(map[string]*runtime.OrderHandle),
	}
}

func (s *OrderHubService) RegisterOrder(req *models.OrderRequest) {
	s.RWmu.Lock()
	defer s.RWmu.Unlock()
	handle := runtime.NewOrderHandle(req)
	s.Orders[req.OrderID] = handle
	go func() {
		err := s.rsOrderService.PushOrder(context.Background(), req)
		if err != nil {
			s.Logger.Error("unable to push order request into the event queue", zap.Error(err))
			return
		}

		s.Logger.Info("order successfully enqueued", zap.String("request", req.OrderID))

		err = handle.OrderFlow.Transition(models.StatusEnqueued)
		if err != nil {
			s.Logger.Error("transition failed", zap.String("request", req.OrderID), zap.Error(err))
			return
		}
	}()
}

func (s *OrderHubService) UnregisterOrder(id string) {
	s.RWmu.Lock()
	defer s.RWmu.Unlock()

	orderHandle, ok := s.Orders[id]

	if !ok {
		s.Logger.Warning("can't find order", zap.String("orderID", id))
		return
	}

	if !orderHandle.OrderFlow.IsTerminated() {
		s.Logger.Fatal("Order delete too early!", zap.String("orderID", id))
		return
	}

	delete(s.Orders, id)

}
