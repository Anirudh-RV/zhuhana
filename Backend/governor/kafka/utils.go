package kafka

import (
	"governor/logger"
	"log"
	"os"
	"strings"
)

var Logger *logger.Logger

func GetKafkaBrokersFromEnv() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		log.Fatalf("No KAFKA_BROKERS set, defaulting to localhost")
		return []string{}
	}
	return strings.Split(brokers, ",")
}

func GetKafkaTopicFromEnv() string {
	return os.Getenv("KAFKA_TOPIC")
}
