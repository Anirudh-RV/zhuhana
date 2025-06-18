package configs

import (
	"os"
	"strings"
)

type KafkaConfig struct {
	Brokers            []string
	OrderRequestTopic  string
	OrderResponseTopic string
	ClientID           string
	GroupID            string
}

func NewKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Brokers:            strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		OrderRequestTopic:  os.Getenv("KAFKA_ORDER_REQUEST_TOPIC"),
		OrderResponseTopic: os.Getenv("KAFKA_ORDER_RESPONSE_TOPIC"),
		ClientID:           os.Getenv("KAFKA_CLIENT_ID"),
		GroupID:            os.Getenv("KAFKA_GROUP_ID"),
	}
}
