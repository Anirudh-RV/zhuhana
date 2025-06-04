package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"governor/logger"
	"log"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

// To create the topic manually. Login to the kafka consumer server and run
// kafka-topics --bootstrap-server governor-kafka:9092 --create \
// --topic cron.jobs --partitions 1 --replication-factor 1

type JobPayload struct {
	JobID   string      `json:"job_id"`
	Target  string      `json:"target"`  // who should consume it, e.g., "governor"
	Payload interface{} `json:"payload"` // job data
	Time    time.Time   `json:"time"`    // time of creation
}

var kafkaClient *kgo.Client

// Initializes the global Kafka client
func InitPublisher(logger *logger.Logger) {
	var err error
	Logger = logger
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
	go Logger.Info(fmt.Sprintf("Publishing CRON job: %s", jobID), zap.String("execution level", "GetAllUserAlgorithms"))
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

func CreateKafkaTopic(topic string, partitions int32, replicationFactor int16) error {
	ctx, cancel := context.WithTimeout(context.Background(), GetKafkaTimeoutFromEnv())
	defer cancel()

	adminClient := kadm.NewClient(kafkaClient)

	// No special configs
	var configs map[string]*string = nil

	// Create the topic
	responses, err := adminClient.CreateTopics(ctx, partitions, replicationFactor, configs, topic)
	if err != nil {
		return fmt.Errorf("admin API call failed: %w", err)
	}

	for _, res := range responses {
		if res.Err != nil {
			if res.Err.Error() == "TOPIC_ALREADY_EXISTS: Topic with this name already exists." {
				go Logger.Info("topic already exists", zap.String("execution level", "CreateKafkaTopic"))
				return nil // Already exists; not an error
			}
			return fmt.Errorf("failed to create topic %s: %w | Error: %s", res.Topic, res.Err, res.Err.Error())
		}
	}

	return nil
}
