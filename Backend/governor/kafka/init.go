package kafka

import (
	"governor/logger"

	"go.uber.org/zap"
)

func Init(logger *logger.Logger) {
	InitLogger(logger)
	go Logger.Info("kafka Logger initialization successful", zap.String("Execution Level", "KafkaInit"))
	InitConsumer()
	go Logger.Info("kafka consumer initialization successful", zap.String("Execution Level", "KafkaInit"))
	InitPublisher()
	go Logger.Info("kafka publisher initialization successful", zap.String("Execution Level", "KafkaInit"))
}
