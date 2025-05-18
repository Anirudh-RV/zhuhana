package services

import (
	"errors"
	"uasam/users/user/models"

	"github.com/alexedwards/argon2id"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (us *UserService) LoginVerifyPasswordHandler(loginVerifyPasswordRequest *models.LoginVerifyPasswordRequest) error {
	passwordHash, err := us.userRepository.GetUserPasswordByEmail(loginVerifyPasswordRequest.EmailID)
	if err != nil {
		go us.logger.Warning("Error accessing user password", zap.String("Execution Level", "LoginVerifyPasswordHandler"), zap.String("Error", err.Error()))
		return err
	}

	passwordMatch, err := argon2id.ComparePasswordAndHash(loginVerifyPasswordRequest.Password, passwordHash)
	if err != nil {
		go us.logger.Warning("Error during password hash comparision", zap.String("Execution Level", "LoginVerifyPasswordHandler"), zap.String("Error", err.Error()))
		return err
	}

	if !passwordMatch {
		go us.logger.Warning("wrong password", zap.String("Execution Level", "LoginVerifyPasswordHandler"))
		return errors.New("wrong password")
	}

	err = us.otpService.SendOTP(loginVerifyPasswordRequest.EmailID)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) LoginVerifyOTPHandler(loginVerifyOTPRequest *models.LoginVerifyOTPRequest) (*models.UserObject, string, error) {
	storedSecret, err := us.otpService.getStoredSecretKey(loginVerifyOTPRequest.EmailID)
	if err == redis.Nil {
		go us.logger.Warning("user secret unavailable in redis", zap.String("execution level", "LoginVerifyOTPHandler"), zap.String("Error", err.Error()))
		return nil, "", err
	} else if err != nil {
		go us.logger.Warning("error in getting user secret in redis", zap.String("execution level", "LoginVerifyOTPHandler"), zap.String("Error", err.Error()))
		return nil, "", err
	}

	status, err := us.otpService.VerifyOTP(storedSecret, loginVerifyOTPRequest.Otp)
	if err != nil {
		go us.logger.Warning("error in verifying otp", zap.String("execution level", "LoginVerifyOTPHandler"), zap.String("Error", err.Error()))
		return nil, "", err
	}
	if !status {
		go us.logger.Warning("wrong otp provided", zap.String("execution level", "LoginVerifyOTPHandler"))
		return nil, "", errors.New("wrong otp provided")
	}

	userObject, err := us.userRepository.GetUserByEmail(loginVerifyOTPRequest.EmailID)
	if err != nil {
		go us.logger.Warning("error in fetching user", zap.String("execution level", "LoginVerifyOTPHandler"), zap.String("Error", err.Error()))
		return nil, "", err
	}

	generatedUserAccessToken, err := us.jwtService.GenerateJWT(userObject.ID.String(), "user")
	if err != nil {
		go us.logger.Warning("errors generating user access token", zap.String("Execution Level", "LoginVerifyOTPHandler"), zap.String("Error", err.Error()))
		return nil, "", err
	}

	return userObject, generatedUserAccessToken, nil
}
