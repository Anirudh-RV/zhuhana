package services

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (uas *UserAlgorithmService) StopUserAlgorithm(userID, userAlgorithmID string) error {
	belongsTo, err := uas.userAlgorthmRepository.DoesUserAlgorithmBelongsToUser(userID, userAlgorithmID)
	if err != nil {
		go uas.logger.Error("could not check ownership of user_algorithm", zap.String("execution level", "StopUserAlgorithm"), zap.String("Error", err.Error()))
		return err
	}
	if !belongsTo {
		go uas.logger.Error("user_algorithm_id does not belong to user_id", zap.String("execution level", "StopUserAlgorithm"))
	}
	userAlgorithmUUID, _ := uuid.Parse(userAlgorithmID)
	go uas.kubernetesService.Stop(userAlgorithmUUID)
	return nil
}
