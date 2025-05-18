package services

import (
	"uasam/users/user/models"

	"github.com/alexedwards/argon2id"
	"go.uber.org/zap"
)

func (us *UserService) ResetPasswordInitHandler(resetPasswordInitRequest *models.ResetPasswordInitRequest) error {
	err := us.otpService.SendResetPasswordEmail(resetPasswordInitRequest.EmailID)
	return err
}

func (us *UserService) ResetPasswordHandler(resetPasswordRequest *models.ResetPasswordRequest) error {
	hashedPassword, err := argon2id.CreateHash(resetPasswordRequest.Password, argon2id.DefaultParams)
	if err != nil {
		go us.logger.Warning("password hashing failed", zap.String("execution level", "ResetPasswordHandler"))
		return err
	}
	err = us.otpService.ResetPassword(resetPasswordRequest.EmailID, resetPasswordRequest.Token, hashedPassword)
	return err
}
