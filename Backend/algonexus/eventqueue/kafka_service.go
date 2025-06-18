package eventqueue

import (
	"algonexus/eventqueue/configs"
	"algonexus/eventqueue/consumers"
	"algonexus/eventqueue/producers"
	"algonexus/logger"
)

type KafkaService struct {
	Producer *producers.KafkaProducer
	Consumer *consumers.KafkaConsumer
	Logger   *logger.Logger
}

func InitKafkaService(callback func([]byte)) (*KafkaService, error) {
	cfg := configs.NewKafkaConfig()
	p, err := producers.NewKafkaProducer(cfg.Brokers, cfg.OrderRequestTopic, cfg.ClientID)
	if err != nil {
		return nil, err
	}

	c, err := consumers.NewOrderConsumer(cfg.Brokers, cfg.OrderResponseTopic, cfg.ClientID)

	return &KafkaService{
		Producer: p,
	}, nil
}
