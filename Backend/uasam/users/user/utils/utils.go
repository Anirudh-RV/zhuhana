package utils

import "uasam/users/user/models"

func MapUserToUserReturnObject(user *models.User) models.UserReturnObject {
	return models.UserReturnObject{
		ID:         user.ID,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		EmailID:    user.EmailID,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}
