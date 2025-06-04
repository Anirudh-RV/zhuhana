package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

// To create the topic manually. Login to the kafka consumer server and run
// kafka-topics --bootstrap-server governor-kafka:9092 --create \
//  --topic cron.jobs --partitions 1 --replication-factor 1

type JobPayload struct {
	JobID   string      `json:"job_id"`
	Target  string      `json:"target"`  // who should consume it, e.g., "governor"
	Payload interface{} `json:"payload"` // job data
	Time    time.Time   `json:"time"`    // time of creation
}

var kafkaClient *kgo.Client

// Initializes the global Kafka client
func InitPublisher() {
	var err error
	brokers := GetKafkaBrokersFromEnv()
	KAFKA_TIMEOUT, _ := strconv.Atoi(os.Getenv("KAFKA_TIMEOUT"))
	kafkaClient, err = kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ProduceRequestTimeout(time.Duration(KAFKA_TIMEOUT)*time.Second),
	)
	if err != nil {
		log.Fatalf("failed to create Kafka client: %v", err)
	}
}

// Publishes a job to a Kafka topic
func PublishJob(jobID string) error {
	go Logger.Info(fmt.Sprintf("Publishing CRON job: %s", jobID), zap.String("execution level", "GetAllUserAlgorithms"))
	job := JobPayload{
		JobID:  jobID,
		Target: "governor",
		Payload: map[string]string{
			"type": "healthcheck",
			"url":  os.Getenv("GOVERNOR_URL"),
		},
		Time: time.Now().UTC(),
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return kafkaClient.ProduceSync(ctx, record).FirstErr()
}
