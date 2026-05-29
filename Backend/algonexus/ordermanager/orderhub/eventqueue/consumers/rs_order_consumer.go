package consumers

import (
	"algonexus/constants"
	"algonexus/logger"
	"algonexus/ordermanager/orderhub/eventqueue"
	"algonexus/ordermanager/orderhub/ports"
	"algonexus/ordermanager/orderhub/registry"
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RsOrderConsumer struct {
	Logger         *logger.Logger
	EventQueue     *eventqueue.RedisStreamEventQueue
	Group          string
	ConsumerName   string
	StreamKey      string
	MessageHandler StreamMessageHandler
	sem            chan struct{}
}

type StreamMessageHandler interface {
	Handle(ctx context.Context, msg redis.XMessage) error
}

func NewRsOrderConsumer(logger *logger.Logger, eventQueue *eventqueue.RedisStreamEventQueue, registry *registry.OrderHubRegistry, broker ports.Broker, stream, group, consumer string) *RsOrderConsumer {
	ctx := context.Background()

	err := eventQueue.Client.XGroupCreateMkStream(ctx, stream, group, "0").Err()
	if err != nil && !redis.HasErrorPrefix(err, "BUSYGROUP") {
		logger.Fatal("Failed to create consumer group",
			zap.String("stream", stream),
			zap.String("group", group),
			zap.String("consumer", consumer),
			zap.String("execution level", "EventQueue Consumer"),
			zap.Error(err),
		)
		panic(err)
	}

	return &RsOrderConsumer{
		Logger:         logger,
		EventQueue:     eventQueue,
		Group:          group,
		ConsumerName:   consumer,
		StreamKey:      stream,
		MessageHandler: NewRsOrderConsumerMsgHandler(logger, registry, broker),
		sem:            make(chan struct{}, constants.AnchorConcurrency),
	}
}

// Run is the polling consumer loop.
func (c *RsOrderConsumer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.pollOnce(ctx, constants.AnchorConcurrency)
		}
	}
}

func (c *RsOrderConsumer) pollOnce(ctx context.Context, count int64) {
	entries, err := c.EventQueue.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    c.Group,
		Consumer: c.ConsumerName,
		Streams:  []string{c.StreamKey, ">"},
		Count:    count,
		Block:    5 * time.Second, // avoid fast polling
	}).Result()

	if err != nil && !errors.Is(err, redis.Nil) {
		c.Logger.Error("XReadGroup failed",
			zap.String("stream", c.StreamKey),
			zap.String("group", c.Group),
			zap.String("consumer", c.ConsumerName),
			zap.String("execution level", "EventQueue Consumer pollOnce"),
			zap.Error(err))
		return
	}

	if entries == nil {
		return
	}

	// Bounded worker pool (AnchorConcurrency): each message is handled in its own
	// goroutine, but the semaphore caps in-flight work. No per-batch WaitGroup barrier,
	// so one slow order never head-of-line blocks the others; when the pool is full the
	// loop stalls on `sem <-`, the next XReadGroup is deferred, and orderstream applies
	// backpressure.
	for _, stream := range entries {
		for _, msg := range stream.Messages {
			c.sem <- struct{}{}
			go func(m redis.XMessage) {
				defer func() { <-c.sem }()
				if err := c.MessageHandler.Handle(ctx, m); err != nil {
					c.Logger.Error("Fail to handle message", zap.String("msgID", m.ID), zap.Error(err))
				}
				if _, err := c.EventQueue.Client.XAck(ctx, c.StreamKey, c.Group, m.ID).Result(); err != nil {
					c.Logger.Error("XAck failed", zap.String("msgID", m.ID), zap.Error(err))
				}
			}(msg)
		}
	}
}
