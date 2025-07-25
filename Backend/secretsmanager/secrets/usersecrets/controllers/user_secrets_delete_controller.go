package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"secretsmanager/logger"
	"secretsmanager/secrets/usersecrets/models"
	"secretsmanager/secrets/usersecrets/services"

	"github.com/gin-gonic/gin"
)

type UserSecretsDeleteController struct {
	userSecretsServiceObj *services.UserSecretsService
	log                   *logger.Logger
}

func NewUserSecretsDeleteController(userSecretsServiceObj *services.UserSecretsService, log *logger.Logger) *UserSecretsDeleteController {
	return &UserSecretsDeleteController{
		userSecretsServiceObj: userSecretsServiceObj,
		log:                   log,
	}
}

func (usc *UserSecretsDeleteController) UserSecretDeleteHandler(c *gin.Context) {
	var userSecretDeleteRequest models.UserSecretsDeleteRequest
	userID, _ := c.Get("USER_ID")
	if userID == nil {
		c.JSON(http.StatusBadRequest, models.UserSecretDeleteResponse{
			Status:            -2,
			StatusDescription: "Unable to parse UserID",
		})
		return
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&userSecretDeleteRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.UserSecretsSetResponse{
			Status:            -1,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	err := usc.userSecretsServiceObj.DeleteUserSecretByID(fmt.Sprint(userID), userSecretDeleteRequest.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.UserSecretDeleteResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, models.UserSecretDeleteResponse{
		Status:            1,
		StatusDescription: "User Secret Deleted",
	})
}
