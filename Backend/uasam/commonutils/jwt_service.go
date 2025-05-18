package commonutils

import (
	"os"
	"strconv"
	"uasam/logger"
)

type JWTService struct {
	logger                            *logger.Logger
	USER_JWT_EXPIRATION_DAYS          int
	USER_JWT_SECRET_KEY               string
	MICROSERVICES_JWT_EXPIRATION_DAYS int
	MICROSERVICES_JWT_SECRET_KEY      string
}

func NewJWTService(logger *logger.Logger) *JWTService {
	USER_JWT_EXPIRATION_DAYS, _ := strconv.Atoi(os.Getenv("USER_JWT_EXPIRATION_DAYS"))
	USER_JWT_SECRET_KEY := os.Getenv("USER_JWT_SECRET_KEY")

	return &JWTService{
		logger:                   logger,
		USER_JWT_EXPIRATION_DAYS: USER_JWT_EXPIRATION_DAYS,
		USER_JWT_SECRET_KEY:      USER_JWT_SECRET_KEY,
	}
}
