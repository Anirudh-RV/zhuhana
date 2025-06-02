package services

import (
	"uasam/users/user/models"

	"go.uber.org/zap"
)

func (us *UserService) AuthenticateUser(userToken string) (*models.UserObject, error) {
	userID, err := us.jwtService.ParseUserJWT(userToken)
	if err != nil {
		go us.logger.Warning("error authenticating jwt token", zap.String("execution level", "AuthenticateUserService"), zap.String("Error", err.Error()))
		return nil, err
	}

	userObjPtr, err := us.userRepository.GetUserByUserID(userID)
	if err != nil {
		go us.logger.Warning("could not get user from id", zap.String("execution level", "AuthenticateUserService"), zap.String("Error", err.Error()))
		return nil, err
	}

	return userObjPtr, nil
}
