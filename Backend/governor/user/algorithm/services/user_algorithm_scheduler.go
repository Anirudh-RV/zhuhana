package services

import (
	"fmt"
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
	uas.schedulerService.ScheduleCronJob(userAlgorithmUUID, startCronSchedule, scheduler.START_USER_ALGORITHM_JOB, uas.kafkaService.GetKafkaTopicFromEnv())
	uas.schedulerService.ScheduleCronJob(userAlgorithmUUID, endCronSchedule, scheduler.END_USER_ALGORITHM_JOB, uas.kafkaService.GetKafkaTopicFromEnv())

	return nil
}

func (uas *UserAlgorithmService) CancelAlgorithmSchedule(userID, userAlgorithmID string) error {
	belongsTo, err := uas.userAlgorthmRepository.DoesUserAlgorithmBelongsToUser(userID, userAlgorithmID)
	if err != nil {
		go uas.logger.Error("could not check ownership of user_algorithm", zap.String("execution level", "CancelAlgorithmSchedule"), zap.String("Error", err.Error()))
		return err
	}
	if !belongsTo {
		go uas.logger.Error("user_algorithm_id does not belong to user_id", zap.String("execution level", "CancelAlgorithmSchedule"))
		return fmt.Errorf("user_algorithm_id does not belong to user_id")
	}
	userAlgorithmUUID, _ := uuid.Parse(userAlgorithmID)
	uas.schedulerService.CancelCronJobForUserAlgorithm(userAlgorithmUUID)
	return nil
}
