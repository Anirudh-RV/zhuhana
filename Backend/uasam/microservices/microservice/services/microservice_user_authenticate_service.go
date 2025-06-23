package services

import (
	"go.uber.org/zap"
)

func (mss *MicroServiceService) AuthenticateMicroServiceUserService(UserServiceToken string) (string, error) {
	userID, err := mss.jwtService.ParseMicroServicesUserJWT(UserServiceToken)
	if err != nil {
		go mss.logger.Warning("error authenticating jwt token", zap.String("execution level", "AuthenticateMicroServiceUserService"), zap.String("Error", err.Error()))
		return "", err
	}

	return userID, nil
}
