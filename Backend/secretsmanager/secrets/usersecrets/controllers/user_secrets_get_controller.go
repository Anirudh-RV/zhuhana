package controllers

import (
	"fmt"
	"net/http"
	"secretsmanager/logger"
	"secretsmanager/secrets/usersecrets/models"
	"secretsmanager/secrets/usersecrets/services"

	"github.com/gin-gonic/gin"
)

type UserSecretsGetController struct {
	userSecretsServiceObj *services.UserSecretsService
	log                   *logger.Logger
}

func NewUserSecretsGetController(userSecretsServiceObj *services.UserSecretsService, log *logger.Logger) *UserSecretsGetController {
	return &UserSecretsGetController{
		userSecretsServiceObj: userSecretsServiceObj,
		log:                   log,
	}
}

// UserSecretsGetHandler godoc
// @Summary      Get a user's secret
// @Description  Fetches a specific secret key-value pair for an authenticated user.
// @Tags         UserSecrets
// @Accept       json
// @Produce      json
// @Param        key query string true "Secret key to fetch"
// @Success      201 {object} models.UserSecretsGetResponse "User secret fetched successfully"
// @Failure      400 {object} models.UserSecretsGetResponse "Invalid request or missing user ID"
// @Failure      500 {object} models.UserSecretsGetResponse "Internal server error"
// @Security     ApiKeyAuth
// @Router       /v1/user/secrets/ [get]
func (usc *UserSecretsSetController) UserSecretsGetHandler(c *gin.Context) {
	var userSecretGetRequest models.UserSecretsGetRequest
	userID, _ := c.Get("USER_ID")
	if userID == nil {
		c.JSON(http.StatusBadRequest, models.UserSecretsGetResponse{
			Status:            -2,
			StatusDescription: "Unable to parse UserID",
		})
		return
	}

	if err := c.ShouldBindQuery(&userSecretGetRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.UserSecretsGetResponse{
			Status:            -1,
			StatusDescription: "Invalid request query parameters",
		})
		return
	}

	userSecret, err := usc.userSecretsServiceObj.GetUserSecret(fmt.Sprint(userID), userSecretGetRequest.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.UserSecretsGetResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}

	c.JSON(http.StatusCreated, models.UserSecretsGetResponse{
		Status:            1,
		StatusDescription: "User Secret Fetched",
		UserSecret:        userSecret,
	})
}
