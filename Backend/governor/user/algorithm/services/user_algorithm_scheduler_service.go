package services

import (
	"fmt"
	"governor/kafka"
	"governor/scheduler"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (uas *UserAlgorithmService) UpdateAlgorithmSchedule(userID, userAlgorithmID, startCronSchedule, endCronSchedule string) error {
	if startCronSchedule == "" || endCronSchedule == "" {
		return fmt.Errorf("cronschedule is empty")
	}

	// TODO: Setup something that will use the CRON to run the container
	err := uas.CronValidator(startCronSchedule)
	if err != nil {
		go uas.logger.Error("could not validate cron schedule / invalid control updated", zap.String("execution level", "UpdateAlgorithmSchedule"), zap.String("Error", err.Error()))
		return err
	}
	err = uas.CronValidator(endCronSchedule)
	if err != nil {
		go uas.logger.Error("could not validate cron schedule / invalid control updated", zap.String("execution level", "UpdateAlgorithmSchedule"), zap.String("Error", err.Error()))
		return err
	}

	err = uas.userAlgorthmRepository.UpdateCronSchedule(userID, userAlgorithmID, startCronSchedule, endCronSchedule)
	if err != nil {
		go uas.logger.Error("could not update algorithm schedule", zap.String("execution level", "UpdateAlgorithmSchedule"), zap.String("Error", err.Error()))
		return err
	}

	userAlgorithmUUID, _ := uuid.Parse(userAlgorithmID)
	scheduler.ScheduleCronJob(userAlgorithmUUID, startCronSchedule, scheduler.START_USER_ALGORITHM_JOB, kafka.GetKafkaTopicFromEnv())
	scheduler.ScheduleCronJob(userAlgorithmUUID, endCronSchedule, scheduler.END_USER_ALGORITHM_JOB, kafka.GetKafkaTopicFromEnv())

	return nil
}
