package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"governor/logger"
	"log"
	"os"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

var cancelConsumer context.CancelFunc

func (kfs *KafkaService) InitConsumer() {
	brokers := kfs.GetKafkaBrokersFromEnv()
	groupID := os.Getenv("KAFKA_GROUP_ID")
	topic := kfs.GetKafkaTopicFromEnv()

	// Start Kafka consumer in background goroutine
	kfs.StartConsumer(brokers, groupID, topic, func(event EventPayload) error {
		kfs.KafkaConsumer(event)
		return nil
	})
}

// Starts the Kafka consumer in a background goroutine
func (kfs *KafkaService) StartConsumer(brokers []string, groupID, topic string, handler func(EventPayload) error) {
	ctx, cancel := context.WithCancel(context.Background())
	cancelConsumer = cancel

	go func() {
		if err := kfs.consumeJobs(ctx, brokers, groupID, topic, handler, kfs.logger); err != nil {
			log.Printf("Kafka consumer exited with error: %v", err)
		}
	}()
}

// Stops the Kafka consumer by cancelling its context
func (kfs *KafkaService) StopConsumer() {
	if cancelConsumer != nil {
		cancelConsumer()
	}
}

func (kfs *KafkaService) consumeJobs(
	ctx context.Context,
	brokers []string,
	groupID string,
	topic string,
	handler func(EventPayload) error,
	logger *logger.Logger,
) error {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		logger.Error("failed to create Kafka consumer client",
			zap.Error(err),
			zap.String("ExecutionLevel", "KafkaInit"),
		)
		return fmt.Errorf("create kafka consumer: %w", err)
	}
	defer client.Close()

	logger.Info("Kafka consumer started",
		zap.String("ExecutionLevel", "KafkaConsumer"),
	)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Kafka consumer context cancelled, shutting down",
				zap.String("ExecutionLevel", "KafkaConsumerShutdown"),
			)
			return nil

		default:
			fetches := client.PollFetches(ctx)
			if fetches.IsClientClosed() {
				return nil
			}

			fetches.EachPartition(func(p kgo.FetchTopicPartition) {
				for _, record := range p.Records {
					var event EventPayload
					if err := json.Unmarshal(record.Value, &event); err != nil {
						logger.Warning("failed to unmarshal Kafka job payload",
							zap.ByteString("rawValue", record.Value),
							zap.Error(err),
							zap.String("ExecutionLevel", "KafkaUnmarshal"),
						)
						continue
					}

					logger.Info("Kafka job received",
						zap.Any("job", event),
						zap.String("ExecutionLevel", "KafkaJobReceived"),
					)

					if err := handler(event); err != nil {
						logger.Error("error in job handler",
							zap.Any("job", event),
							zap.Error(err),
							zap.String("ExecutionLevel", "KafkaJobHandler"),
						)
					}
				}
			})

			if err := client.CommitMarkedOffsets(ctx); err != nil {
				logger.Warning("failed to commit Kafka offsets",
					zap.Error(err),
					zap.String("ExecutionLevel", "KafkaOffsetCommit"),
				)
			}

			client.AllowRebalance()
		}
	}
}
