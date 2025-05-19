package services

import (
	"uasam/commonutils"
	"uasam/logger"
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
