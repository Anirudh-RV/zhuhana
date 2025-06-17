package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kadm"
	"go.uber.org/zap"
)

func (kfs *KafkaService) CreateKafkaTopic(topic string, partitions int32, replicationFactor int16) error {
	ctx, cancel := context.WithTimeout(context.Background(), kfs.GetKafkaTimeoutFromEnv())
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
				go kfs.logger.Info("topic already exists", zap.String("execution level", "CreateKafkaTopic"))
				return nil // Already exists; not an error
			}
			return fmt.Errorf("failed to create topic %s: %w | Error: %s", res.Topic, res.Err, res.Err.Error())
		}
	}

	return nil
}
