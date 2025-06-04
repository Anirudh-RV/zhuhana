package kafka

import (
	"governor/logger"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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

func GetKafkaTimeoutFromEnv() time.Duration {
	KAFKA_TIMEOUT, _ := strconv.Atoi(os.Getenv("KAFKA_TIMEOUT"))
	return time.Duration(KAFKA_TIMEOUT) * time.Second
}
