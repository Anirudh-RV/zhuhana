package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (uas *UserAlgorithmService) StartUserAlgorithm(userID, userAlgorithmID, market, symbol string, startTime, endTime *time.Time, portfolioSize, frequency int) error {
	belongsTo, err := uas.userAlgorthmRepository.DoesUserAlgorithmBelongsToUser(userID, userAlgorithmID)
	if err != nil {
		go uas.logger.Error("could not check ownership of user_algorithm", zap.String("execution level", "StartUserAlgorithm"), zap.String("Error", err.Error()))
		return err
	}
	if !belongsTo {
		go uas.logger.Error("user_algorithm_id does not belong to user_id", zap.String("execution level", "StartUserAlgorithm"))
		return fmt.Errorf("user_algorithm_id does not belong to user_id")
	}
	userAlgorithmUUID, _ := uuid.Parse(userAlgorithmID)
	go uas.kubernetesService.Start(userAlgorithmUUID, market, symbol, startTime, endTime, portfolioSize, frequency)
	return nil
}
