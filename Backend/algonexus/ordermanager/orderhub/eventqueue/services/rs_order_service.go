package services

import (
	"algonexus/constants"
	"algonexus/logger"
	"algonexus/ordermanager/models"
	"algonexus/ordermanager/orderhub/eventqueue"
	"algonexus/ordermanager/orderhub/eventqueue/consumers"
	"algonexus/ordermanager/orderhub/eventqueue/producers"
	"algonexus/ordermanager/orderhub/ports"
	"algonexus/ordermanager/orderhub/registry"
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

// RsOrderService owns the OrderHub ingress: the producer writes orders to the single
// OrderStream, and the anchor consumer (a bounded submit-worker pool) drives the FSM to
// SUBMITTED and hands each order to the broker via the in-process Broker port.
type RsOrderService struct {
	registry   *registry.OrderHubRegistry
	logger     *logger.Logger
	eventQueue *eventqueue.RedisStreamEventQueue
	producer   *producers.RsOrderProducer
	consumer   *consumers.RsOrderConsumer
}

func NewRsOrderService(logger *logger.Logger, registry *registry.OrderHubRegistry, broker ports.Broker) *RsOrderService {
	ctx := context.Background()
	rsEq := eventqueue.NewRedisStreamEventQueue(ctx, logger)

	p := producers.NewRsOrderProducer(logger, rsEq)

	c := consumers.NewRsOrderConsumer(
		logger, rsEq, registry, broker,
		constants.OrderStream, constants.OrderStreamGroup, constants.OrderStreamConsumer,
	)

	return &RsOrderService{
		registry:   registry,
		logger:     logger,
		eventQueue: rsEq,
		producer:   p,
		consumer:   c,
	}
}

func (s *RsOrderService) StartAll(ctx context.Context) {
	go func() {
		s.logger.Info("Starting order anchor consumer", zap.String("stream", constants.OrderStream))
		s.consumer.Run(ctx)
	}()
}

func (s *RsOrderService) PushOrderNonWait(ctx context.Context, request *models.OrderRequest) error {
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		s.logger.Error("json marshal failed", zap.Error(err))
		return fmt.Errorf("json marshal failed: %w", err)
	}
	values := map[string]interface{}{
		"data": string(jsonBytes),
	}

	if err := s.producer.Produce(ctx, constants.OrderStream, values); err != nil {
		s.logger.Error("produce failed", zap.Error(err))
		return fmt.Errorf("produce failed: %w", err)
	}
	return nil
}
