package services

import (
	"go.uber.org/zap"
)

func (mss *MicroServiceService) AuthenticateUserMicroService(userToken string) (string, error) {
	userID, err := mss.jwtService.ParseMicroServicesUserJWT(userToken)
	if err != nil {
		go mss.logger.Warning("error authenticating jwt token", zap.String("execution level", "AuthenticateUserMicroService"), zap.String("Error", err.Error()))
		return "", err
	}

	return userID, nil
}
