package services

import (
	"go.uber.org/zap"
)

func (uss *UserSecretsService) DeleteUserSecretByID(userID, secretID string) error {
	err := uss.userSecretsRepository.DeleteUserSecretByID(userID, secretID)
	if err != nil {
		go uss.logger.Warning("error deleting user secret", zap.String("execution level", "DeleteUserSecretByID"), zap.String("Error", err.Error()))
		return err
	}
	return err
}
