package producers

import (
	"algonexus/logger"
	"algonexus/ordermanager/orderhub/eventqueue"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Redis Stream Producer
type RsOrderProducer struct {
	EventQueue *eventqueue.RedisStreamEventQueue
	Logger     *logger.Logger
}

func NewRsOrderProducer(logger *logger.Logger, rsEventQueue *eventqueue.RedisStreamEventQueue) *RsOrderProducer {
	return &RsOrderProducer{
		EventQueue: rsEventQueue,
		Logger:     logger,
	}
}

func (p *RsOrderProducer) Produce(ctx context.Context, stream string, values map[string]interface{}) error {
	args := &redis.XAddArgs{
		Stream: stream,
		Values: values,
	}

	id, err := p.EventQueue.Client.XAdd(ctx, args).Result()
	if err != nil {
		p.Logger.Error("XAdd Failed",
			zap.String("stream", stream),
			zap.String("execution level", "EventQueue Producer"),
			zap.Error(err),
		)
		return fmt.Errorf("XAdd failed: %w", err)
	}

	p.Logger.Info("Producer XAdd Done",
		zap.String("stream", stream),
		zap.String("msgID", id),
		zap.String("execution level", "EventQueue Producer"),
	)

	return nil
}
