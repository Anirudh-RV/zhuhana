package services

import (
	"encoding/json"
	"time"
	"uasam/users/user/models"
	"uasam/users/user/utils"

	"github.com/alexedwards/argon2id"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (us *UserService) SignUpInitHandler(signUpRequestObject *models.SignUpInitRequest) error {
	err := us.otpService.SendOTP(signUpRequestObject.EmailID)
	if err != nil {
		return err
	}

	hashedPassword, err := argon2id.CreateHash(signUpRequestObject.Password, argon2id.DefaultParams)
	if err != nil {
		go us.logger.Warning("password hashing failed", zap.String("execution level", "signUpInitHandler"))
		return err
	}
	signUpRequestObject.Password = hashedPassword

	jsonSignUpRequestObject, err := json.Marshal(signUpRequestObject)
	if err != nil {
		return err
	}

	err = us.redis.Set(*us.ctx, signUpRequestObject.EmailID, jsonSignUpRequestObject, time.Duration(us.otpService.OTP_DURATION_IN_SECONDS)*time.Second).Err()
	if err != nil {
		go us.logger.Warning("could not store user object in redis", zap.String("execution level", "signUpInitHandler"))
		return err
	}

	return nil
}

func (us *UserService) SignUpVerifyOTPHandler(signUpVerifyOTPRequest *models.SignUpVerifyOTPRequest) (*models.UserObject, string, int, error) {
	storedSecret, err := us.otpService.getStoredSecretKey(signUpVerifyOTPRequest.EmailID)
	if err == redis.Nil {
		go us.logger.Warning("user secret unavailable in redis", zap.String("execution level", "SignUpVerifyOTPHandler"), zap.String("Error", err.Error()))
		return nil, "", 0, err
	} else if err != nil {
		go us.logger.Warning("error in getting user secret in redis", zap.String("execution level", "SignUpVerifyOTPHandler"), zap.String("Error", err.Error()))
		return nil, "", 0, err
	}

	status, err := us.otpService.VerifyOTP(storedSecret, signUpVerifyOTPRequest.Otp)
	if err != nil {
		go us.logger.Warning("error in verifying otp", zap.String("execution level", "SignUpVerifyOTPHandler"), zap.String("Error", err.Error()))
		return nil, "", 0, err
	}
	if !status {
		go us.logger.Warning("wrong otp provided", zap.String("execution level", "SignUpVerifyOTPHandler"))
		return nil, "", -1, err
	}

	userJSON, err := us.redis.Get(*us.ctx, signUpVerifyOTPRequest.EmailID).Result()
	if err != nil {
		go us.logger.Warning("could not get user object from redis", zap.String("Execution Level", "SignUpVerifyOTPHandler"), zap.String("Error", err.Error()))
		return nil, "", 0, err
	}

	var signUpInitRequestObject models.SignUpInitRequest
	err = json.Unmarshal([]byte(userJSON), &signUpInitRequestObject)
	if err != nil {
		go us.logger.Warning("json decoding error for user object", zap.String("Execution Level", "SignUpVerifyOTPHandler"), zap.String("Error", err.Error()))
		return nil, "", 0, err
	}

	userObject, err := us.userRepository.CreateUser(signUpInitRequestObject.FirstName, signUpInitRequestObject.MiddleName, signUpInitRequestObject.LastName, signUpInitRequestObject.EmailID, signUpInitRequestObject.Password)
	if err != nil {
		go us.logger.Warning("errors creating user object", zap.String("Execution Level", "SignUpVerifyOTPHandler"), zap.String("Error", err.Error()))
		return nil, "", 0, err
	}

	userResponseObject := utils.MapUserToUserObject(userObject)

	generatedUserAccessToken, err := us.jwtService.GenerateUserJWT(userObject.ID.String(), "user")
	if err != nil {
		go us.logger.Warning("errors generating user access token", zap.String("Execution Level", "SignUpVerifyOTPHandler"), zap.String("Error", err.Error()))
		return nil, "", 0, err
	}

	return userResponseObject, generatedUserAccessToken, 1, nil
}
