package services

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (uas *UserAlgorithmService) StartUserAlgorithm(userID, userAlgorithmID string) error {
	belongsTo, err := uas.userAlgorthmRepository.DoesUserAlgorithmBelongsToUser(userID, userAlgorithmID)
	if err != nil {
		go uas.logger.Error("could not check ownership of user_algorithm", zap.String("execution level", "StartUserAlgorithm"), zap.String("Error", err.Error()))
		return err
	}
	if !belongsTo {
		go uas.logger.Error("user_algorithm_id does not belong to user_id", zap.String("execution level", "StartUserAlgorithm"))
	}
	userAlgorithmUUID, _ := uuid.Parse(userAlgorithmID)
	go uas.kubernetesService.Start(userAlgorithmUUID)
	return nil
}
