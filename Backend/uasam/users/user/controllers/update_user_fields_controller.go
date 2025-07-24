package controllers

import (
	"encoding/json"
	"net/http"
	"uasam/logger"
	"uasam/users/user/models"
	userService "uasam/users/user/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateUserFieldsController struct {
	userService *userService.UserService
	log         *logger.Logger
}

func NewUpdateUserFieldsController(userService *userService.UserService, log *logger.Logger) *UpdateUserFieldsController {
	return &UpdateUserFieldsController{
		userService: userService,
		log:         log,
	}
}

func (ufc *UpdateUserFieldsController) UpdateUserFieldsHandler(c *gin.Context) {
	var updateUserFieldsRequest models.UpdateUserFieldsRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&updateUserFieldsRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.UpdateUserFieldsResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	rawUserID, _ := c.Get("USER_ID")
	userIDStr, ok := rawUserID.(string)
	if !ok {
		c.JSON(http.StatusOK, &models.UpdateUserFieldsResponse{
			Status:            -1,
			StatusDescription: "unable to find USER_ID",
		})
		return
	}

	// Parse to UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid USER_ID format"})
		return
	}

	err = ufc.userService.UpdateUserNameFields(userID, updateUserFieldsRequest.FirstName, updateUserFieldsRequest.MiddleName, updateUserFieldsRequest.LastName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.UpdateUserFieldsResponse{
			Status:            0,
			StatusDescription: "User fields update error",
		})
		return
	}

	c.JSON(http.StatusOK, &models.UpdateUserFieldsResponse{
		Status:            1,
		StatusDescription: "User fields updated",
	})
}
