package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

var kafkaClient *kgo.Client

// Initializes the global Kafka client
func (kfs *KafkaService) InitPublisher() {
	var err error
	brokers := kfs.GetKafkaBrokersFromEnv()
	kafkaClient, err = kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ProduceRequestTimeout(kfs.GetKafkaTimeoutFromEnv()),
	)
	if err != nil {
		log.Fatalf("failed to create Kafka client: %v", err)
	}

	err = kfs.CreateKafkaTopic(kfs.GetKafkaTopicFromEnv(), 1, 1)
	if err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}
}

// Publishes a job to a Kafka topic
func (kfs *KafkaService) PublishJob(eventID, eventType string, payload interface{}) error {
	go kfs.logger.Info(fmt.Sprintf("publishing cron job: %s", eventID), zap.String("execution level", "GetAllUserAlgorithms"))
	event := EventPayload{
		EventID:   eventID,
		EventType: eventType,
		Target:    os.Getenv("ORIGIN_SERVICE"),
		Payload:   payload,
		Time:      time.Now().UTC(),
	}

	val, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Send the message
	record := &kgo.Record{
		Topic: kfs.GetKafkaTopicFromEnv(),
		Key:   []byte(eventID),
		Value: val,
	}

	ctx, cancel := context.WithTimeout(context.Background(), kfs.GetKafkaTimeoutFromEnv())
	defer cancel()

	return kafkaClient.ProduceSync(ctx, record).FirstErr()
}
