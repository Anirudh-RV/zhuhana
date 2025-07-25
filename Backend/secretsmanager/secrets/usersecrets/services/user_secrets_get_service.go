package services

import (
	"secretsmanager/secrets/usersecrets/models"

	"go.uber.org/zap"
)

func (uss *UserSecretsService) GetUserSecret(userID, key string) (*models.UserSecret, error) {

	userSecret, err := uss.userSecretsRepository.GetUserSecret(userID, key)
	if err != nil {
		go uss.logger.Warning("error getting user secret", zap.String("execution level", "GetUserSecret"), zap.String("Error", err.Error()))
		return nil, err
	}

	if userSecret == nil {
		go uss.logger.Warning("no key exists", zap.String("execution level", "GetUserSecret"))
		return nil, nil
	}

	return userSecret, nil
}

func (uss *UserSecretsService) GetUserKeys(userID string) ([]string, error) {
	userKeys, err := uss.userSecretsRepository.GetAllUserSecretKeys(userID)
	if err != nil {
		go uss.logger.Warning("error getting user keys", zap.String("execution level", "GetUserKeys"), zap.String("Error", err.Error()))
		return nil, err
	}

	return userKeys, nil
}
