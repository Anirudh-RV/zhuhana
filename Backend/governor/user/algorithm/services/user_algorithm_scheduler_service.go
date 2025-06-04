package services

import (
	"fmt"

	"go.uber.org/zap"
)

func (uas *UserAlgorithmService) UpdateAlgorithmSchedule(userID, scriptID, cronSchedule string) error {
	if cronSchedule == "" {
		return fmt.Errorf("cronschedule is empty")
	}

	// TODO: Setup something that will use the CRON to run the container
	err := uas.CronValidator(cronSchedule)
	if err != nil {
		go uas.logger.Error("could not validate cron schedule / invalid control updated", zap.String("execution level", "UpdateAlgorithmSchedule"), zap.String("Error", err.Error()))
		return err
	}
	err = uas.userAlgorthmRepository.UpdateCronSchedule(userID, scriptID, cronSchedule)
	if err != nil {
		go uas.logger.Error("could not update algorithm schedule", zap.String("execution level", "UpdateAlgorithmSchedule"), zap.String("Error", err.Error()))
		return err
	}

	return nil
}
