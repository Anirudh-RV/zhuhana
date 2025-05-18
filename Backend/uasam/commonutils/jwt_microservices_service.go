package commonutils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func (jts *JWTService) GenerateMicroServicesJWT(callerMicroService string, calledMicroService string) (string, error) {
	claims := jwt.MapClaims{
		"caller_microservice": callerMicroService,
		"called_microservice": calledMicroService,
		"exp":                 time.Now().Add(time.Duration(jts.MICROSERVICES_JWT_EXPIRATION_DAYS) * 24 * time.Hour).Unix(),
		"iat":                 time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jts.MICROSERVICES_JWT_SECRET_KEY))
}

func (jts *JWTService) ParseMicroServicesJWT(tokenStr string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jts.MICROSERVICES_JWT_SECRET_KEY, nil
	})

	if err != nil || !token.Valid {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		go jts.logger.Warning("claims for otp not correct", zap.String("execution level", "ParseMicroServicesJWT"))
		return "", "", jwt.ErrInvalidKey
	}

	callerMicroService, ok := claims["caller_microservice"].(string)
	if !ok {
		go jts.logger.Warning("caller_microservice retrieval from jwt token failed", zap.String("execution level", "ParseMicroServicesJWT"))
		return "", "", jwt.ErrInvalidKey
	}

	calledMicroService, ok := claims["called_microservice"].(string)
	if !ok {
		go jts.logger.Warning("called_microservice retrieval from jwt token failed", zap.String("execution level", "ParseMicroServicesJWT"))
		return "", "", jwt.ErrInvalidKey
	}

	return callerMicroService, calledMicroService, nil
}
