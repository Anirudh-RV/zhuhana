package services

import (
	"errors"
	"uasam/commonutils"
	"uasam/logger"

	microServicesConstants "uasam/microservices/constants"

	"go.uber.org/zap"
)

type MicroServiceService struct {
	logger     *logger.Logger
	jwtService *commonutils.JWTService
}

func NewMicroServiceService(logger *logger.Logger, jwtService *commonutils.JWTService) *MicroServiceService {

	return &MicroServiceService{
		logger:     logger,
		jwtService: jwtService,
	}
}

func (mss *MicroServiceService) GenerateMicroServiceAccessKey(callerService string, apiKey string) (string, error) {
	calleeService, err := mss.jwtService.CheckMicroServiceAPIKey(apiKey)
	if err != nil {
		return "", err
	}

	_, exists := microServicesConstants.ALL_SERVICES[callerService]
	if !exists {
		go mss.logger.Warning("invalid caller_service", zap.String("execution level", "GenerateMicroServiceAccessKey"))
		return "", errors.New("invalid caller_service")
	}

	accessToken, err := mss.jwtService.GenerateMicroServicesJWT(callerService, calleeService)
	if err != nil {
		go mss.logger.Warning("error generating jwt token", zap.String("execution level", "GenerateMicroServiceAccessKey"))
		return "", err
	}

	return accessToken, nil
}
