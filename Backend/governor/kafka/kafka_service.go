package kafka

import (
	"governor/kubernetescontroller"
	"governor/logger"

	"go.uber.org/zap"
)

func (kfs *KafkaService) Init(logger *logger.Logger) {
	go kfs.logger.Info("kafka Logger initialization successful", zap.String("Execution Level", "KafkaInit"))
	kfs.InitConsumer()
	go kfs.logger.Info("kafka consumer initialization successful", zap.String("Execution Level", "KafkaInit"))
	kfs.InitPublisher()
	go kfs.logger.Info("kafka publisher initialization successful", zap.String("Execution Level", "KafkaInit"))
}

type KafkaService struct {
	logger            *logger.Logger
	kubernetesService *kubernetescontroller.KubernetesService
}

func NewKafkaService(logger *logger.Logger, kubernetesService *kubernetescontroller.KubernetesService) *KafkaService {
	return &KafkaService{
		logger:            logger,
		kubernetesService: kubernetesService,
	}
}
