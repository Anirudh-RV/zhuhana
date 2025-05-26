package services

import (
	"go.uber.org/zap"
)

func (mss *MicroServiceService) GenerateMicroServiceUserAccessKey(userID string) (string, error) {
	accessToken, err := mss.jwtService.GenerateMicroServicesUserJWT(userID)
	if err != nil {
		go mss.logger.Warning("error generating jwt token", zap.String("execution level", "GenerateMicroServiceUserAccessKey"), zap.String("Error", err.Error()))
		return "", err
	}

	return accessToken, nil
}
