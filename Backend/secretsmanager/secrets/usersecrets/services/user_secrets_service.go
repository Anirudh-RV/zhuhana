package services

import (
	"secretsmanager/logger"
	"secretsmanager/secrets/usersecrets/repositories"
)

type UserSecretsService struct {
	logger                *logger.Logger
	userSecretsRepository *repositories.UserSecretRepository
}

func NewUserSecretsService(logger *logger.Logger, userSecretsRepository *repositories.UserSecretRepository) *UserSecretsService {

	return &UserSecretsService{
		logger:                logger,
		userSecretsRepository: userSecretsRepository,
	}
}
