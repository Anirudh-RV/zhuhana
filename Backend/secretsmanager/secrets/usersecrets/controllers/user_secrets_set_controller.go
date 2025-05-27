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

// UserSecretsSetHandler godoc
// @Summary      Set a user's secret
// @Description  Stores or updates a secret key-value pair for an authenticated user.
// @Tags         UserSecrets
// @Accept       json
// @Produce      json
// @Param        userSecret body models.UserSecretsSetRequest true "Secret key-value to set"
// @Success      201 {object} models.UserSecretsSetResponse "User secret set successfully"
// @Failure      400 {object} models.UserSecretsSetResponse "Invalid request or missing user ID"
// @Failure      500 {object} models.UserSecretsSetResponse "Internal server error"
// @Security     ApiKeyAuth
// @Router       /v1/user/secrets/ [post]
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
		c.JSON(http.StatusBadRequest, models.UserSecretsSetResponse{
			Status:            -1,
			StatusDescription: "Invalid request payload",
		})
		return
	}

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
