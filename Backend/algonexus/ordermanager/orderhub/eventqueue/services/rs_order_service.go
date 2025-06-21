package services

import (
	"algonexus/logger"
	"algonexus/ordermanager/models"
	"algonexus/ordermanager/orderhub/eventqueue"
	"algonexus/ordermanager/orderhub/eventqueue/consumers"
	"algonexus/ordermanager/orderhub/eventqueue/producers"
	"algonexus/ordermanager/orderhub/registry"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

type RsOrderService struct {
	strategyName string
	registry     *registry.OrderHubRegistry
	logger       *logger.Logger

	consumer *consumers.RsOrderConsumer // Single consumer
	producer *producers.RsOrderProducer
}

func NewRsOrderService(logger *logger.Logger, registry *registry.OrderHubRegistry) *RsOrderService {
	ctx := context.Background()
	rsEq := eventqueue.NewRedisStreamEventQueue(ctx, logger)

	//Centric producer
	p := producers.NewRsOrderProducer(logger, rsEq)

	// Only one consumer for now
	var strategy = "strategy-1"
	var streamKey = fmt.Sprintf("orderstream:%s", strategy)
	var groupName = fmt.Sprintf("group:%s", strategy)
	var consumerName = fmt.Sprintf("consumer-%s", strategy)

	c := consumers.NewRsOrderConsumer(logger, rsEq, registry, streamKey, groupName, consumerName)

	return &RsOrderService{
		strategyName: strategy,
		consumer:     c,
		producer:     p,
		logger:       logger,
	}
}

func (s *RsOrderService) StartAll(ctx context.Context) {
	go func(str string, c *consumers.RsOrderConsumer) {
		s.logger.Info("Starting consumer", zap.String("strategy", str))
		c.Run(ctx)
	}("strategy-1", s.consumer)
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

	stream := fmt.Sprintf("orderstream:%s", s.strategyName)

	s.logger.Info("order is going to be produced", zap.String("msg", string(jsonBytes)))

	err = s.producer.Produce(ctx, stream, values)
	if err != nil {
		s.logger.Error("produce failed", zap.Error(err))
		return fmt.Errorf("produce failed: %w", err)
	}
	return nil
}
