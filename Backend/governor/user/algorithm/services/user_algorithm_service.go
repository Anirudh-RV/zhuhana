package services

import (
	"fmt"
	"governor/commonutils"
	"governor/logger"
	"governor/middleware"
	"governor/user/algorithm/models"
	"governor/user/algorithm/repositories"

	"go.uber.org/zap"
)

type UserAlgorithmService struct {
	logger                    *logger.Logger
	userAlgorthmRepository    *repositories.UserAlgorithmRepository
	microserviceAuthenticator *middleware.MicroSeviceAuthenticator
}

func NewUserAlgorithmService(logger *logger.Logger, userAlgorthmRepository *repositories.UserAlgorithmRepository, microserviceAuthenticator *middleware.MicroSeviceAuthenticator) *UserAlgorithmService {
	return &UserAlgorithmService{
		logger:                    logger,
		userAlgorthmRepository:    userAlgorthmRepository,
		microserviceAuthenticator: microserviceAuthenticator,
	}
}

func (uas *UserAlgorithmService) CronValidator(cronSchedule string) error {
	// TODO: Validate the timing of the cron somehow
	_, err := commonutils.ValidateAndParseCron(cronSchedule)
	if err != nil {
		return err
	}
	return nil
}

func (uas *UserAlgorithmService) GetAllUserAlgorithms(userID string) ([]models.UserAlgorithmInfo, error) {
	userAlgorithms, err := uas.userAlgorthmRepository.GetAllUserAlgorithmByUserID(userID)
	if err != nil {
		go uas.logger.Error("error getting user algorithm", zap.String("execution level", "GetAllUserAlgorithms"), zap.String("Error", err.Error()))
		return nil, err
	}

	for _, userAlgorithm := range userAlgorithms {
		if userAlgorithm.ScriptURL != nil {
			presignedURL, err := commonutils.GetPresignedURL(userID, fmt.Sprint(userAlgorithm.ScriptID))
			if err != nil {
				go uas.logger.Error("error getting presigned url for script", zap.String("execution level", "GetAllUserAlgorithms"), zap.String("Error", err.Error()))
				return nil, err
			}
			userAlgorithm.ScriptURL = &presignedURL
		}
	}

	return userAlgorithms, nil
}

func (uas *UserAlgorithmService) GetUserAlgorithm(userID, algorithmID string) (*models.UserAlgorithmInfo, error) {
	userAlgorithm, err := uas.userAlgorthmRepository.GetUserAlgorithmByUserID(userID, algorithmID)
	if err != nil {
		go uas.logger.Error("error getting user algorithm", zap.String("execution level", "GetUserAlgorithm"), zap.String("Error", err.Error()))
		return nil, err
	}

	if userAlgorithm.ScriptURL != nil {
		presignedURL, err := commonutils.GetPresignedURL(userID, fmt.Sprint(userAlgorithm.ScriptID))
		if err != nil {
			go uas.logger.Error("error getting presigned url for script", zap.String("execution level", "GetUserAlgorithm"), zap.String("Error", err.Error()))
			return nil, err
		}
		userAlgorithm.ScriptURL = &presignedURL
	}

	return userAlgorithm, nil
}
