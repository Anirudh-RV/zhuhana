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

// UpdateUserFieldsHandler godoc
//
// @Summary      Update user name fields
// @Description  Updates first name, middle name, and last name of the authenticated user
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.UpdateUserFieldsRequest true "User fields to update"
// @Success      200 {object} models.UpdateUserFieldsResponse "User fields updated successfully"
// @Failure      400 {object} models.UpdateUserFieldsResponse "Invalid request payload or USER_ID format"
// @Failure      401 {object} models.UpdateUserFieldsResponse "Unauthorized - update failed"
// @Router       /v1/user/edit/ [put]
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
		c.JSON(http.StatusOK, &models.UpdateUserFieldsResponse{
			Status:            -1,
			StatusDescription: "unable to find USER_ID",
		})
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
