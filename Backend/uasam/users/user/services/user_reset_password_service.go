package services

import (
	"fmt"
	"uasam/users/user/models"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (us *UserService) ResetPasswordInitHandler(resetPasswordInitRequest *models.ResetPasswordInitRequest, device, ipAddress string) error {
	err := us.otpService.SendResetPasswordEmail(resetPasswordInitRequest.EmailID, device, ipAddress)
	return err
}

func (us *UserService) ResetPasswordHandler(resetPasswordRequest *models.ResetPasswordRequest) error {
	if resetPasswordRequest.Password == "" {
		go us.logger.Warning("password cannot be empty", zap.String("execution level", "ResetPasswordHandler"))
		return fmt.Errorf("password cannot be empty")
	}
	hashedPassword, err := argon2id.CreateHash(resetPasswordRequest.Password, argon2id.DefaultParams)
	if err != nil {
		go us.logger.Warning("password hashing failed", zap.String("execution level", "ResetPasswordHandler"))
		return err
	}
	err = us.otpService.ResetPassword(resetPasswordRequest.EmailID, resetPasswordRequest.Token, hashedPassword)
	return err
}

func (us *UserService) UpdatePasswordHandler(id uuid.UUID, password string) error {
	if password == "" {
		go us.logger.Warning("password cannot be empty", zap.String("execution level", "UpdatePasswordHandler"))
		return fmt.Errorf("password cannot be empty")
	}

	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		go us.logger.Warning("password hashing failed", zap.String("execution level", "UpdatePasswordHandler"))
		return err
	}

	err = us.userRepository.UpdateUserPasswordByID(id, hashedPassword)
	if err != nil {
		go us.logger.Warning("update password failed", zap.String("execution level", "UpdatePasswordHandler"))
		return err
	}

	return err
}
