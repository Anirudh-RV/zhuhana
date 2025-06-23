package services

import (
	"go.uber.org/zap"
)

func (mss *MicroServiceService) GenerateMicroServiceUserAlgorithmAccessKey(userAlgorithmID string) (string, error) {
	accessToken, err := mss.jwtService.GenerateMicroServicesUserAlgorithmJWT(userAlgorithmID)
	if err != nil {
		go mss.logger.Warning("error generating jwt token", zap.String("execution level", "GenerateMicroServiceUserAlgorithmAccessKey"), zap.String("Error", err.Error()))
		return "", err
	}

	return accessToken, nil
}
