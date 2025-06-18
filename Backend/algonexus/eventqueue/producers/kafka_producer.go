package producers

import (
	"context"
	"encoding/json"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaProducer struct {
	Client *kgo.Client
	Topic  string
}

func NewKafkaProducer(brokers []string, topic string, clientID string) (*KafkaProducer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ProducerLinger(10e6), // batching
		kgo.ClientID(clientID),
	)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		Client: client,
		Topic:  topic,
	}, nil
}

func (p *KafkaProducer) SendOrder(ctx context.Context, key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	record := &kgo.Record{
		Topic: p.Topic,
		Key:   []byte(key),
		Value: b,
	}

	return p.Client.ProduceSync(ctx, record).FirstErr()
}
