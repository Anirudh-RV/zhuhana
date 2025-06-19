package eventqueue

import (
	"algonexus/eventqueue/consumers"
	"algonexus/eventqueue/producers"
	"algonexus/logger"
	"context"
	"fmt"
	"go.uber.org/zap"
)

type RsOrderService struct {
	consumer *consumers.RsOrderConsumer // Single consumer
	producer *producers.RsOrderProducer
	logger   *logger.Logger
}

func NewRsOrderService(logger *logger.Logger) *RsOrderService {
	ctx := context.Background()
	rsEq := NewRedisStreamEventQueue(ctx, logger)

	//Centric producer
	p := producers.NewRsOrderProducer(logger, rsEq)

	// Only one consumer for now
	var strategy = "strategy-1"
	var streamKey = fmt.Sprintf("orderstream:%s", strategy)
	var groupName = fmt.Sprintf("group:%s", strategy)
	var consumerName = fmt.Sprintf("consumer-%s", strategy)

	c := consumers.NewRsOrderConsumer(logger, rsEq, streamKey, groupName, consumerName)

	return &RsOrderService{
		consumer: c,
		producer: p,
		logger:   logger,
	}
}

func (s *RsOrderService) StartAll(ctx context.Context) {
	go func(str string, c *consumers.RsOrderConsumer) {
		s.logger.Info("Starting consumer", zap.String("strategy", str))
		c.Run(ctx)
	}("strategy-1", s.consumer)
}
