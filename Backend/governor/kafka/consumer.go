package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"go.uber.org/zap"
)

func (kfs *KafkaService) KafkaConsumer(event EventPayload) error {
	// TODO: Get the kubernetes to run the container
	log.Printf("kafka consumer received event %+v", event)
	switch event.EventType {
	case CRON_JOB_EVENT_TYPE:
		jsonData, err := json.Marshal(event.Payload)
		if err != nil {
			fmt.Println("marshal error:", err)
			return err
		}

		// Unmarshal into CronJob
		var cronJob CronJob
		if err := json.Unmarshal(jsonData, &cronJob); err != nil {
			fmt.Println("unmarshal error:", err)
			return err
		}
		kfs.CronJobConsumer(cronJob)

	default:
		go kfs.logger.Info("event type not recognized", zap.String("execution level", "KafkaConsumer"))

	}
	return nil
}

func (kfs *KafkaService) CronJobConsumer(cronJob CronJob) {
	go kfs.logger.Info(fmt.Sprintf("processing cronJob %+v", cronJob), zap.String("execution level", "KafkaConsumer"))
	switch cronJob.JobType {
	case START_USER_ALGORITHM_JOB:
		kfs.kubernetesService.CronStart(cronJob.UserAlgorithmID)
	case END_USER_ALGORITHM_JOB:
		kfs.kubernetesService.Stop(cronJob.UserAlgorithmID)
	default:
		go kfs.logger.Info("job_type not recognized", zap.String("execution level", "CronJobConsumer"))
	}
}
