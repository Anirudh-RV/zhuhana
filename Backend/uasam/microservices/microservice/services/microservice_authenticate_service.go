package services

import (
	"errors"

	"go.uber.org/zap"
)

func (mss *MicroServiceService) AuthenticateMicroService(originService string, AuthToken string) (string, string, error) {
	callerMicroService, calledMicroService, err := mss.jwtService.ParseMicroServicesJWT(AuthToken)
	if err != nil {
		go mss.logger.Warning("error authenticating jwt token", zap.String("execution level", "AuthenticateMicroService"), zap.String("Error", err.Error()))
		return "", "", err
	}

	if callerMicroService != originService {
		go mss.logger.Warning("not authorized for this service", zap.String("execution level", "AuthenticateMicroService"))
		return "", "", errors.New("not authorized for this service")
	}

	return callerMicroService, calledMicroService, nil
}
