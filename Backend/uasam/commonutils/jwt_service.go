package commonutils

import (
	"os"
	"strconv"
	"time"
	"uasam/logger"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type JWTService struct {
	logger              *logger.Logger
	JWT_EXPIRATION_DAYS int
	JWT_SECRET_KEY      string
}

func NewJWTService(logger *logger.Logger) *JWTService {
	JWT_EXPIRATION_DAYS, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_DAYS"))
	JWT_SECRET_KEY := os.Getenv("JWT_SECRET_KEY")

	return &JWTService{
		logger:              logger,
		JWT_EXPIRATION_DAYS: JWT_EXPIRATION_DAYS,
		JWT_SECRET_KEY:      JWT_SECRET_KEY,
	}
}

func (jts *JWTService) GenerateJWT(userID string, userType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_type": userType,
		"exp":       time.Now().Add(time.Duration(jts.JWT_EXPIRATION_DAYS) * 24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jts.JWT_SECRET_KEY))
}

func (jts *JWTService) ParseJWT(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jts.JWT_SECRET_KEY, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		go jts.logger.Warning("claims for otp not correct", zap.String("Execution Level", "ParseJWT"))
		return "", jwt.ErrInvalidKey
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		go jts.logger.Warning("user_id retrieval from jwt token failed", zap.String("Execution Level", "ParseJWT"))
		return "", jwt.ErrInvalidKey
	}

	return userID, nil
}
