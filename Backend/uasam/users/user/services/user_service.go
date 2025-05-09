package services

import (
	"context"
	"encoding/json"
	"time"
	"uasam/users/user/models"
	"uasam/users/user/repositories"

	"uasam/logger"

	"github.com/alexedwards/argon2id"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserService struct {
	ctx            *context.Context
	otpService     *OTPService
	userRepository *repositories.UserRepository
	logger         *logger.Logger
	redis          *redis.Client
}

func NewUserService(ctx *context.Context, otpService *OTPService, userRepository *repositories.UserRepository, logger *logger.Logger, redis *redis.Client) *UserService {

	return &UserService{
		ctx:            ctx,
		otpService:     otpService,
		userRepository: userRepository,
		logger:         logger,
		redis:          redis,
	}
}

func (us *UserService) CreateUser(firstname string, middleName string, lastName string, emailID string, password string) (*models.User, error) {
	var middleNamePtr *string
	if middleName != "" {
		middleNamePtr = &middleName
	}

	user, err := us.userRepository.CreateUser(firstname, middleNamePtr, lastName, emailID, password)

	if err != nil {

		return nil, err
	}

	return user, nil
}

func (us *UserService) IfUserExists(emailID string) (bool, error) {
	status, err := us.userRepository.IfUserEmailExists(emailID)
	if err != nil {
		return false, err
	}
	return status, nil
}

func (us *UserService) SignUpInitHandler(signUpRequestObject *models.SignUpInitRequest) error {
	err := us.otpService.SendOTP(signUpRequestObject.EmailID)
	if err != nil {
		return err
	}

	hashedPassword, err := argon2id.CreateHash(signUpRequestObject.Password, argon2id.DefaultParams)
	if err != nil {
		go us.logger.Warning("Password hashing failed", zap.String("Execution Level", "signUpInitHandler"))
		return err
	}
	signUpRequestObject.Password = hashedPassword

	jsonSignUpRequestObject, err := json.Marshal(signUpRequestObject)
	if err != nil {
		return err
	}

	err = us.redis.Set(*us.ctx, signUpRequestObject.EmailID, jsonSignUpRequestObject, time.Duration(us.otpService.OTP_DURATION_IN_SECONDS)*time.Second).Err()
	if err != nil {
		go us.logger.Warning("Could not store User Object in Redis", zap.String("Execution Level", "signUpInitHandler"))
		return err
	}

	return nil
}
