package services

import (
	"algonexus/constants"
	"algonexus/logger"
	brokermodels "algonexus/ordermanager/backtestengine/broker/models"
	"algonexus/ordermanager/models"
	eqServices "algonexus/ordermanager/orderhub/eventqueue/services"
	"algonexus/ordermanager/orderhub/ports"
	"algonexus/ordermanager/orderhub/registry"
	"algonexus/ordermanager/orderhub/runtime"
	"context"
	"sync"

	"go.uber.org/zap"
)

type OrderHubService struct {
	logger    *logger.Logger
	registry  *registry.OrderHubRegistry
	eqService *eqServices.RsOrderService
	broker    ports.Broker
	rwmu      sync.RWMutex
}

func NewOrderHubService(logger *logger.Logger, broker ports.Broker) *OrderHubService {
	hubRegistry := registry.NewOrderHubRegistry(logger)

	rsOrderService := eqServices.NewRsOrderService(logger, hubRegistry, broker)
	logger.Info("RedisStreams event queue init successful", zap.String("Execution Level", "OrderHub"))

	rsOrderService.StartAll(context.Background())
	logger.Info("RedisStreams event queue is running", zap.String("Execution Level", "OrderHub"))

	s := &OrderHubService{
		logger:    logger,
		registry:  hubRegistry,
		eqService: rsOrderService,
		broker:    broker,
	}

	s.startListeners(constants.FillsConcurrency)
	logger.Info("OrderHub Listener pool started",
		zap.Int("workers", constants.FillsConcurrency),
		zap.String("Execution Level", "OrderHub"))

	return s
}

// startListeners launches the shared Listener pool: FillsConcurrency goroutines all
// draining the broker's single fill channel. Whichever worker grabs a fill finalizes
// that order (terminal FSM + registry.Delete). There is no per-order goroutine, so the
// goroutine count is fixed regardless of in-flight order volume (no leak).
func (s *OrderHubService) startListeners(n int) {
	for i := 0; i < n; i++ {
		go s.listen()
	}
}

func (s *OrderHubService) listen() {
	for evt := range s.broker.Fills() {
		s.finalize(evt)
	}
}

// finalize drives the terminal FSM tail based on the broker's verdict and removes the
// handle from the pool. It is the single place that completes an order and cleans up.
func (s *OrderHubService) finalize(evt brokermodels.BrokerFillEvent) {
	handle := s.registry.Get(evt.OrderID)
	if handle == nil {
		// Already finalized/cleaned, or a fill for a fail-fast-rejected order — skip.
		s.logger.Warning("fill for unknown/closed order", zap.String("orderID", evt.OrderID))
		return
	}

	if evt.Status == models.StatusBrokerFilled {
		if err := handle.OrderFlow.Transition(models.StatusBrokerConfirmed); err != nil {
			s.logger.Error("transition failed", zap.String("orderID", evt.OrderID), zap.String("to", string(models.StatusBrokerConfirmed)), zap.Error(err))
		}
		// Happy path assumes a single full fill -> COMPLETED. Partial fills deferred.
		if err := handle.OrderFlow.Transition(models.StatusComplete); err != nil {
			s.logger.Error("transition failed", zap.String("orderID", evt.OrderID), zap.String("to", string(models.StatusComplete)), zap.Error(err))
		}
		s.logger.Info("order complete", zap.String("orderID", evt.OrderID))
	} else {
		if err := handle.OrderFlow.Transition(models.StatusError); err != nil {
			s.logger.Error("transition failed", zap.String("orderID", evt.OrderID), zap.String("to", string(models.StatusError)), zap.Error(err))
		}
		s.logger.Warning("broker reported non-fill", zap.String("orderID", evt.OrderID), zap.String("status", string(evt.Status)))
	}

	s.registry.Delete(evt.OrderID)
}

func (s *OrderHubService) RegisterOrder(req *models.OrderRequest) {
	handle := runtime.NewOrderHandle(req)
	s.registry.Update(req.OrderID, handle)

	// Producer goroutine. ENQUEUED is recorded BEFORE the XADD: once the order is on
	// OrderStream the anchor can drive ASSIGNED concurrently, so publishing first would
	// let the anchor race ahead of ENQUEUED.
	go func() {
		if err := handle.OrderFlow.Transition(models.StatusPendingSend); err != nil {
			s.logger.Error("transition failed", zap.String("request", req.OrderID), zap.Error(err))
			return
		}
		if err := handle.OrderFlow.Transition(models.StatusEnqueued); err != nil {
			s.logger.Error("transition failed", zap.String("request", req.OrderID), zap.Error(err))
			return
		}
		if err := s.eqService.PushOrderNonWait(context.Background(), req); err != nil {
			s.logger.Error("unable to push order request into the event queue", zap.Error(err))
			return
		}
		s.logger.Info("order published to OrderStream", zap.String("request", req.OrderID))
	}()
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
