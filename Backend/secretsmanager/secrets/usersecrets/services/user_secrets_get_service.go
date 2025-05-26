package services

import (
	"secretsmanager/secrets/usersecrets/models"

	"go.uber.org/zap"
)

func (uss *UserSecretsService) GetUserSecret(userID, key, value string) (*models.UserSecret, error) {

	userSecret, err := uss.userSecretsRepository.GetUserSecret(userID, key)
	if err != nil {
		go uss.logger.Warning("error getting user secret", zap.String("execution level", "SetUserSecrets"), zap.String("Error", err.Error()))
		return nil, err
	}

	return userSecret, nil
}
