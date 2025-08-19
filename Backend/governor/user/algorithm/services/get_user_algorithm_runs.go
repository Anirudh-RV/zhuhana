package services

import (
	"fmt"
	"governor/user/algorithm/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (uas *UserAlgorithmService) GetUserAlgorithmRuns(userID, userAlgorithmID string) ([]models.UserAlgorithmRun, error) {
	belongsTo, err := uas.userAlgorthmRepository.DoesUserAlgorithmBelongsToUser(userID, userAlgorithmID)
	if err != nil {
		go uas.logger.Error("could not check ownership of user_algorithm", zap.String("execution level", "StartUserAlgorithm"), zap.String("Error", err.Error()))
		return nil, err
	}
	if !belongsTo {
		go uas.logger.Error("user_algorithm_id does not belong to user_id", zap.String("execution level", "StartUserAlgorithm"))
		return nil, fmt.Errorf("user_algorithm_id does not belong to user_id")
	}
	userAlgorithmUUID, _ := uuid.Parse(userAlgorithmID)
	userAlgorithmRuns, err := uas.userAlgorthmRepository.GetUserAlgorithmRunsByUserAlgorithmID(userAlgorithmUUID)
	if err != nil {
		go uas.logger.Error("user_algorithm_id does not belong to user_id", zap.String("execution level", "StartUserAlgorithm"))
		return nil, err
	}
	return userAlgorithmRuns, nil
}
