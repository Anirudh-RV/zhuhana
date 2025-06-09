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
func InitPublisher() {
	var err error
	brokers := GetKafkaBrokersFromEnv()
	kafkaClient, err = kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ProduceRequestTimeout(GetKafkaTimeoutFromEnv()),
	)
	if err != nil {
		log.Fatalf("failed to create Kafka client: %v", err)
	}

	err = CreateKafkaTopic(GetKafkaTopicFromEnv(), 1, 1)
	if err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}
}

// Publishes a job to a Kafka topic
func PublishJob(jobID string, payload interface{}) error {
	go Logger.Info(fmt.Sprintf("publishing cron job: %s", jobID), zap.String("execution level", "GetAllUserAlgorithms"))
	job := JobPayload{
		JobID:   jobID,
		Target:  os.Getenv("ORIGIN_SERVICE"),
		Payload: payload,
		Time:    time.Now().UTC(),
	}

	val, err := json.Marshal(job)
	if err != nil {
		return err
	}

	// Send the message
	record := &kgo.Record{
		Topic: GetKafkaTopicFromEnv(),
		Key:   []byte(jobID),
		Value: val,
	}

	ctx, cancel := context.WithTimeout(context.Background(), GetKafkaTimeoutFromEnv())
	defer cancel()

	return kafkaClient.ProduceSync(ctx, record).FirstErr()
}
