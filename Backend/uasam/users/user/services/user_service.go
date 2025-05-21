package services

import (
	"context"
	"uasam/commonutils"
	"uasam/users/user/models"
	"uasam/users/user/repositories"

	"uasam/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserService struct {
	ctx            *context.Context
	otpService     *OTPService
	jwtService     *commonutils.JWTService
	userRepository *repositories.UserRepository
	logger         *logger.Logger
	redis          *redis.Client
}

func NewUserService(ctx *context.Context, otpService *OTPService, jwtService *commonutils.JWTService, userRepository *repositories.UserRepository, logger *logger.Logger, redis *redis.Client) *UserService {

	return &UserService{
		ctx:            ctx,
		otpService:     otpService,
		jwtService:     jwtService,
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
		go us.logger.Warning("error checking if user exists", zap.String("execution level", "IfUserExists"), zap.String("Error", err.Error()))
		return false, err
	}
	return status, nil
}
