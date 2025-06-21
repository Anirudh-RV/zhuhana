package services

import (
	"algonexus/logger"
	"algonexus/ordermanager/models"
	eqServices "algonexus/ordermanager/orderhub/eventqueue/services"
	"algonexus/ordermanager/orderhub/registry"
	"algonexus/ordermanager/orderhub/runtime"
	"context"
	"go.uber.org/zap"
	"sync"
)

type OrderHubService struct {
	logger    *logger.Logger
	registry  *registry.OrderHubRegistry
	eqService *eqServices.RsOrderService
	rwmu      sync.RWMutex
}

func NewOrderHubService(logger *logger.Logger) *OrderHubService {
	hubRegistry := registry.NewOrderHubRegistry(logger)

	rsOrderService := eqServices.NewRsOrderService(logger, hubRegistry)
	logger.Info("RedisStreams event queue init successful", zap.String("Execution Level", "OrderHub"))

	rsOrderService.StartAll(context.Background())
	logger.Info("RedisStreams event queue is running", zap.String("Execution Level", "OrderHub"))

	return &OrderHubService{
		logger:    logger,
		registry:  hubRegistry,
		eqService: rsOrderService,
	}
}

func (s *OrderHubService) Listen(id string) {
	handle := s.registry.Get(id)

	if handle == nil {
		s.logger.Warning("OrderHub attempted to listen to unknown order", zap.String("id", id))
		return
	}

	for event := range handle.Channel {
		//TODO Stop Condition
		s.logger.Info("OrderHub received event", zap.String("order_id", id), zap.String("type", string(event.Type)))
	}
}

func (s *OrderHubService) RegisterOrder(req *models.OrderRequest) {
	orderSession := runtime.NewOrderSession(req)
	s.registry.Update(req.OrderID, orderSession)

	// Producer
	go func() {
		err := orderSession.OrderFlow.Transition(models.StatusPendingSend)
		if err != nil {
			s.logger.Error("transition failed", zap.String("request", req.OrderID), zap.Error(err))
			return
		}

		err = s.eqService.PushOrderNonWait(context.Background(), req)
		if err != nil {
			s.logger.Error("unable to push order request into the event queue", zap.Error(err))
			return
		}

		s.logger.Info("order successfully enqueued", zap.String("request", req.OrderID))

		err = orderSession.OrderFlow.Transition(models.StatusEnqueued)
		if err != nil {
			s.logger.Error("transition failed", zap.String("request", req.OrderID), zap.Error(err))
			return
		}
	}()

	// Listener routine (listen to consumer)
	// TODO: Refactor listener logic later
	// TODO: Relocate per-order handling logic into OrderHandle for encapsulation
	s.logger.Info("started a listener", zap.String("orderId", req.OrderID))
	go s.Listen(req.OrderID)

}

func (s *OrderHubService) UnregisterOrder(id string) {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	handle := s.registry.Get(id)

	if handle == nil {
		s.logger.Warning("can't find order", zap.String("orderID", id))
		return
	}

	if !handle.OrderFlow.IsTerminated() {
		s.logger.Fatal("Order delete too early!", zap.String("orderID", id))
		return
	}

	s.registry.Delete(id)

}
