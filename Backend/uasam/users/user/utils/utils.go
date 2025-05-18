package utils

import "uasam/users/user/models"

func MapUserToUserObject(user *models.User) *models.UserObject {
	return &models.UserObject{
		ID:         user.ID,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		EmailID:    user.EmailID,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}
