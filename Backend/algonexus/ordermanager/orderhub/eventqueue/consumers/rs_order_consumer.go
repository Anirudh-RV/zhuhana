package consumers

import (
	"algonexus/logger"
	"algonexus/ordermanager/orderhub/eventqueue"
	"algonexus/ordermanager/orderhub/registry"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"sync"
	"time"
)

type RsOrderConsumer struct {
	Logger         *logger.Logger
	EventQueue     *eventqueue.RedisStreamEventQueue
	Registry       *registry.OrderHubRegistry
	Group          string
	ConsumerName   string
	StreamKey      string
	MessageHandler StreamMessageHandler
	WaitGroup      *sync.WaitGroup
}

type StreamMessageHandler interface {
	Handle(ctx context.Context, msg redis.XMessage) error
}

func NewRsOrderConsumer(logger *logger.Logger, eventQueue *eventqueue.RedisStreamEventQueue, registry *registry.OrderHubRegistry, stream, group, consumer string) *RsOrderConsumer {
	ctx := context.Background()

	err := eventQueue.Client.XGroupCreateMkStream(ctx, stream, group, "0").Err() // Sole consumer

	// Non-BUSYGROUP error
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
		WaitGroup:      &sync.WaitGroup{},
		MessageHandler: NewRsOrderConsumerMsgHandler(logger, registry),
	}
}

// Run Polling Consumer
func (c *RsOrderConsumer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.pollOnce(ctx, 5)
		}
	}
}

func (c *RsOrderConsumer) pollOnce(ctx context.Context, count int64) {
	entries, err := c.EventQueue.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    c.Group,
		Consumer: c.ConsumerName,
		Streams:  []string{c.StreamKey, ">"},
		Count:    count,           //set to 1 to ensure sequential read
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
		c.Logger.Warning("XReadGroup returned nil entries")
		return
	}

	for _, stream := range entries {
		for _, msg := range stream.Messages {
			c.WaitGroup.Add(1)
			go func(m redis.XMessage) {
				defer c.WaitGroup.Done()
				c.Logger.Info("Message received", zap.String("msgID", msg.ID))
				err := c.MessageHandler.Handle(ctx, msg)
				if err != nil {
					c.Logger.Error("Fail to handle message", zap.String("msgID", msg.ID), zap.Error(err))
				}

				_, err = c.EventQueue.Client.XAck(ctx, c.StreamKey, c.Group, msg.ID).Result()

				if err != nil {
					c.Logger.Error("XAck failed", zap.String("msgID", msg.ID), zap.Error(err))
				}
			}(msg)
		}
		c.WaitGroup.Wait()
	}

}
