package commonutils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func (jts *JWTService) GenerateUserJWT(userID string, userType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_type": userType,
		"exp":       time.Now().Add(time.Duration(jts.USER_JWT_EXPIRATION_DAYS) * 24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jts.USER_JWT_SECRET_KEY))
}

func (jts *JWTService) ParseUserJWT(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jts.USER_JWT_SECRET_KEY), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		go jts.logger.Warning("claims for otp not correct", zap.String("execution level", "ParseUserJWT"))
		return "", jwt.ErrInvalidKey
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		go jts.logger.Warning("user_id retrieval from jwt token failed", zap.String("execution level", "ParseUserJWT"))
		return "", jwt.ErrInvalidKey
	}

	return userID, nil
}
