package services

import (
	"errors"
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
		return nil, errors.New("no key exists")
	}

	return userSecret, nil
}
