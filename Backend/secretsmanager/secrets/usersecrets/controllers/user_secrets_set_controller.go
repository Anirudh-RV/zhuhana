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

type UserSecretsSetController struct {
	userSecretsServiceObj *services.UserSecretsService
	log                   *logger.Logger
}

func NewUserSecretsSetController(userSecretsServiceObj *services.UserSecretsService, log *logger.Logger) *UserSecretsSetController {
	return &UserSecretsSetController{
		userSecretsServiceObj: userSecretsServiceObj,
		log:                   log,
	}
}

func (usc *UserSecretsSetController) UserSecretsSetHandler(c *gin.Context) {
	var userSecretSetRequest models.UserSecretsSetRequest
	userID, _ := c.Get("USER_ID")
	if userID == nil {
		c.JSON(http.StatusBadRequest, models.UserSecretsSetResponse{
			Status:            -2,
			StatusDescription: "Unable to parse UserID",
		})
		return
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&userSecretSetRequest); err != nil {
		fmt.Println("ERROR: ", err.Error())
		c.JSON(http.StatusBadRequest, models.UserSecretsSetResponse{
			Status:            -1,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	fmt.Println("USER ID FROM AUTH: ", fmt.Sprint(userID))

	err := usc.userSecretsServiceObj.SetUserSecret(fmt.Sprint(userID), userSecretSetRequest.Key, userSecretSetRequest.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.UserSecretsSetResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}

	c.JSON(http.StatusCreated, models.UserSecretsSetResponse{
		Status:            1,
		StatusDescription: "User Secret Added",
	})
}
