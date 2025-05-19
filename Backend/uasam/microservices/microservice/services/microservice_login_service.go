package services

import (
	"errors"

	microServicesConstants "uasam/microservices/constants"

	"go.uber.org/zap"
)

func (mss *MicroServiceService) GenerateMicroServiceAccessKey(originService string, apiKey string) (string, error) {
	calleeService, err := mss.jwtService.CheckMicroServiceAPIKey(apiKey)
	if err != nil {
		return "", err
	}

	_, exists := microServicesConstants.ALL_SERVICES[originService]
	if !exists {
		go mss.logger.Warning("invalid caller_service", zap.String("execution level", "GenerateMicroServiceAccessKey"))
		return "", errors.New("invalid caller_service")
	}

	accessToken, err := mss.jwtService.GenerateMicroServicesJWT(originService, calleeService)
	if err != nil {
		go mss.logger.Warning("error generating jwt token", zap.String("execution level", "GenerateMicroServiceAccessKey"), zap.String("Error", err.Error()))
		return "", err
	}

	return accessToken, nil
}
