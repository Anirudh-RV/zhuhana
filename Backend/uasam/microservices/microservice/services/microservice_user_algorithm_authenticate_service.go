package services

import (
	"go.uber.org/zap"
)

func (mss *MicroServiceService) AuthenticateMicroServiceUserAlgorithm(UserAlgorithmToken string) (string, error) {
	userAlgorithmID, err := mss.jwtService.ParseMicroServicesUserAlgorithmJWT(UserAlgorithmToken)
	if err != nil {
		go mss.logger.Warning("error authenticating jwt token", zap.String("execution level", "AuthenticateMicroServiceUserService"), zap.String("Error", err.Error()))
		return "", err
	}

	return userAlgorithmID, nil
}
