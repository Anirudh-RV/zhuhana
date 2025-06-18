package consumers

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaConsumer struct {
	Client *kgo.Client
}

func NewOrderConsumer(brokers []string, topic string, groupID string, clientID string) (*KafkaConsumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumeTopics(topic),
		kgo.ConsumerGroup(groupID),
		kgo.ClientID(clientID+"-consumer"),
		kgo.AutoCommitMarks(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init Kafka consumer: %w", err)
	}

	return &KafkaConsumer{
		Client: client,
	}, nil
}

// StartPolling starts consuming messages from the configured topic
// and invokes the provided handler(msg []byte) on each message.
func (c *KafkaConsumer) StartPolling(ctx context.Context, handler func([]byte)) {
	go func() {
		for {
			fetches := c.Client.PollFetches(ctx)

			// You can handle errors globally if needed
			fetches.EachPartition(func(p kgo.FetchTopicPartition) {
				for _, record := range p.Records {
					handler(record.Value)
					c.Client.MarkCommitRecords(record) // AutoCommit enabled
				}
			})

		}
	}()
}
