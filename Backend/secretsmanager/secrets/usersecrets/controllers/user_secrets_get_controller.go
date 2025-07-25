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
// @Router       /v1/user/secret/ [get]
func (usc *UserSecretsGetController) UserSecretGetHandler(c *gin.Context) {
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

	c.JSON(http.StatusOK, models.UserSecretsGetResponse{
		Status:            1,
		StatusDescription: "User Secret Fetched",
		UserSecret:        userSecret,
	})
}

// UserSecretKeysGetHandler godoc
//
// @Summary      Get all secret keys for the current user
// @Description  Returns all secret keys stored for the authenticated user.
// @Tags         user-secrets
// @Produce      json
// @Success      200  {object}  models.UserSecretKeysResponse  "Success"
// @Failure      400  {object}  models.UserSecretKeysResponse  "Invalid user ID"
// @Failure      500  {object}  models.UserSecretKeysResponse  "Server error"
// @Router       /v1/user/secret/keys [get]
// @Security     ApiKeyAuth
func (usc *UserSecretsGetController) UserSecretKeysGetHandler(c *gin.Context) {
	userID, _ := c.Get("USER_ID")
	if userID == nil {
		c.JSON(http.StatusBadRequest, models.UserSecretKeysResponse{
			Status:            -2,
			StatusDescription: "Unable to parse UserID",
		})
		return
	}

	userSecretKeys, err := usc.userSecretsServiceObj.GetUserKeys(fmt.Sprint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.UserSecretKeysResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, models.UserSecretKeysResponse{
		Status:            1,
		StatusDescription: "User Secret Keys Fetched",
		Keys:              userSecretKeys,
	})
}
