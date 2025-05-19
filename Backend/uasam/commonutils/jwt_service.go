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
	GOVERNOR_API_KEY                  string
	ORCHESTRATOR_API_KEY              string
	OUTBOUND_API_KEY                  string
	UASAM_API_KEY                     string
	SECRETS_MANAGER_API_KEY           string
	ALL_API_KEYS                      map[string]string
}

func NewJWTService(logger *logger.Logger) *JWTService {
	USER_JWT_EXPIRATION_DAYS, _ := strconv.Atoi(os.Getenv("USER_JWT_EXPIRATION_DAYS"))
	MICROSERVICES_JWT_EXPIRATION_DAYS, _ := strconv.Atoi(os.Getenv("MICROSERVICES_JWT_EXPIRATION_DAYS"))
	USER_JWT_SECRET_KEY := os.Getenv("USER_JWT_SECRET_KEY")

	GOVERNOR_API_KEY := os.Getenv("GOVERNOR_API_KEY")
	ORCHESTRATOR_API_KEY := os.Getenv("ORCHESTRATOR_API_KEY")
	OUTBOUND_API_KEY := os.Getenv("OUTBOUND_API_KEY")
	UASAM_API_KEY := os.Getenv("UASAM_API_KEY")
	SECRETS_MANAGER_API_KEY := os.Getenv("SECRETS_MANAGER_API_KEY")
	ALL_API_KEYS := map[string]string{
		GOVERNOR_API_KEY:        "governor",
		ORCHESTRATOR_API_KEY:    "orchestrator",
		OUTBOUND_API_KEY:        "outbound",
		UASAM_API_KEY:           "uasam",
		SECRETS_MANAGER_API_KEY: "secrets-manager",
	}

	return &JWTService{
		logger:                            logger,
		USER_JWT_EXPIRATION_DAYS:          USER_JWT_EXPIRATION_DAYS,
		MICROSERVICES_JWT_EXPIRATION_DAYS: MICROSERVICES_JWT_EXPIRATION_DAYS,
		USER_JWT_SECRET_KEY:               USER_JWT_SECRET_KEY,
		GOVERNOR_API_KEY:                  GOVERNOR_API_KEY,
		ORCHESTRATOR_API_KEY:              ORCHESTRATOR_API_KEY,
		OUTBOUND_API_KEY:                  OUTBOUND_API_KEY,
		UASAM_API_KEY:                     UASAM_API_KEY,
		SECRETS_MANAGER_API_KEY:           SECRETS_MANAGER_API_KEY,
		ALL_API_KEYS:                      ALL_API_KEYS,
	}
}
