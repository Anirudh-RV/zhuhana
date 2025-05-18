package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	"go.uber.org/zap"
)

func (ots *OTPService) generatePasswordResetToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (ots *OTPService) hashPasswordResetToken(token string) (string, error) {
	hashedToken, err := argon2id.CreateHash(token, argon2id.DefaultParams)
	if err != nil {
		go ots.logger.Warning("token hashing failed", zap.String("execution level", "signUpInitHandler"))
		return "", err
	}

	return hashedToken, nil
}

func (ots *OTPService) storeResetPasswordHash(emailID, token string) error {
	err := ots.redis.Set(*ots.ctx, emailID+ots.RESET_PASSWORD_KEY_SUFFIX, token, time.Duration(ots.RESET_PASSWORD_TOKEN_DURATION_IN_SECOND)*time.Second).Err()
	if err != nil {
		go ots.logger.Warning("could not store password reset token in redis", zap.String("execution level", "storeResetPasswordHash"), zap.String("Error", err.Error()))
		return err
	}
	return nil
}

func (ots *OTPService) SendResetPasswordEmail(emailID string) error {
	token, err := ots.generatePasswordResetToken()
	if err != nil {
		return err
	}

	hashedToken, err := ots.hashPasswordResetToken(token)
	if err != nil {
		return err
	}

	err = ots.storeResetPasswordHash(emailID, hashedToken)
	if err != nil {
		return err
	}

	err = ots.emailService.SendResetPasswordEmail(emailID, token)
	if err != nil {
		return err
	}

	return nil
}

func (ots *OTPService) ResetPassword(emailID, token, newPasswordHash string) error {
	retreivedHashedToken, err := ots.redis.Get(*ots.ctx, emailID+ots.RESET_PASSWORD_KEY_SUFFIX).Result()
	if err != nil {
		go ots.logger.Warning("could not get password reset token in redis", zap.String("execution level", "ResetPassword"), zap.String("Error", err.Error()))
		return err
	}

	tokenMatch, err := argon2id.ComparePasswordAndHash(token, retreivedHashedToken)
	if err != nil {
		go ots.logger.Warning("error during token hash comparision", zap.String("execution level", "ResetPassword"), zap.String("Error", err.Error()))
		return err
	}
	if !tokenMatch {
		go ots.logger.Warning("wrong token", zap.String("execution level", "ResetPassword"))
		return errors.New("wrong token")
	}

	err = ots.userRepository.UpdateUserPasswordByEmail(emailID, newPasswordHash)
	if err != nil {
		go ots.logger.Warning("error during password update", zap.String("execution level", "ResetPassword"), zap.String("Error", err.Error()))
		return err
	}

	return nil
}
