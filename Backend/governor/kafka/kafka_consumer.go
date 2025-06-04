package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"governor/logger"
	"log"
	"os"

	"github.com/twmb/franz-go/pkg/kgo"
)

var cancelConsumer context.CancelFunc

func InitConsumer(logger *logger.Logger) {
	brokers := GetKafkaBrokersFromEnv()
	groupID := os.Getenv("KAFKA_GROUP_ID")
	topic := GetKafkaTopicFromEnv()
	Logger = logger

	// Start Kafka consumer in background goroutine
	StartConsumer(brokers, groupID, topic, func(job JobPayload) error {
		log.Printf("Received job %s", job.JobID)
		return nil
	})
}

// Starts the Kafka consumer in a background goroutine
func StartConsumer(brokers []string, groupID, topic string, handler func(JobPayload) error) {
	ctx, cancel := context.WithCancel(context.Background())
	cancelConsumer = cancel

	go func() {
		if err := consumeCronJobs(ctx, brokers, groupID, topic, handler); err != nil {
			log.Printf("Kafka consumer exited with error: %v", err)
		}
	}()
}

// Stops the Kafka consumer by cancelling its context
func StopConsumer() {
	if cancelConsumer != nil {
		cancelConsumer()
	}
}

// The actual consumer logic that continuously polls Kafka
func consumeCronJobs(
	ctx context.Context,
	brokers []string,
	groupID string,
	topic string,
	handler func(JobPayload) error,
) error {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		return fmt.Errorf("failed to create kafka consumer client: %w", err)
	}
	defer client.Close()

	log.Println("[KafkaConsumer] Consumer started...")

	for {
		select {
		case <-ctx.Done():
			log.Println("[KafkaConsumer] Context cancelled, shutting down consumer.")
			return nil
		default:
			fetches := client.PollFetches(ctx)
			if fetches.IsClientClosed() {
				return nil // clean shutdown
			}

			fetches.EachPartition(func(p kgo.FetchTopicPartition) {
				for _, record := range p.Records {
					var job JobPayload
					if err := json.Unmarshal(record.Value, &job); err != nil {
						log.Printf("[KafkaConsumer] Failed to unmarshal: %v", err)
						continue
					}
					log.Printf("[KafkaConsumer] Received job: %+v", job)

					// Call your handler
					if err := handler(job); err != nil {
						log.Printf("[KafkaConsumer] Job handler error: %v", err)
					}
				}
			})

			// Mark all messages as processed
			client.AllowRebalance()
		}
	}
}
