package services

import (
	"governor/logger"
	"governor/user/algorithm/models"
	"governor/user/algorithm/repositories"
	"mime/multipart"
)

type UserAlgorithmService struct {
	logger                 *logger.Logger
	userAlgorthmRepository *repositories.UserRepository
}

func NewUserAlgorithmService(logger *logger.Logger, userAlgorthmRepository *repositories.UserRepository) *UserAlgorithmService {
	return &UserAlgorithmService{
		logger:                 logger,
		userAlgorthmRepository: userAlgorthmRepository,
	}
}

func (pbs *UserAlgorithmService) UserAlgorithmHandler(userID, ScriptName, cronSchedule string, script multipart.File) (*models.UserAlgorithm, error) {
	// TODO
	return nil, nil
}
