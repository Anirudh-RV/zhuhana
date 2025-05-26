package controllers

import (
	"encoding/json"
	"net/http"
	"uasam/logger"
	"uasam/microservices/microservice/models"
	microServiceService "uasam/microservices/microservice/services"

	"github.com/gin-gonic/gin"
)

type MicroServiceUserLoginController struct {
	microServiceServiceObj *microServiceService.MicroServiceService
	log                    *logger.Logger
}

func NewMicroServiceUserLoginController(microServiceServiceObj *microServiceService.MicroServiceService, log *logger.Logger) *MicroServiceUserLoginController {
	return &MicroServiceUserLoginController{
		microServiceServiceObj: microServiceServiceObj,
		log:                    log,
	}
}

// MicroServiceUserLoginHandler godoc
// @Summary Authenticate a microservice user
// @Description Generates an access token for a valid microservice user
// @Tags Microservice Authentication
// @Accept json
// @Produce json
// @Param user body models.MicroServiceUserLoginRequest true "Microservice User Login Request"
// @Success 200 {object} models.MicroServiceUserLoginResponse
// @Failure 400 {object} models.MicroServiceUserLoginResponse "Invalid request payload"
// @Failure 500 {object} models.MicroServiceUserLoginResponse "Server error"
// @Router /v1/microservice/user/login/ [post]
func (mulc *MicroServiceUserLoginController) MicroServiceUserLoginHandler(c *gin.Context) {
	var userLoginRequest models.MicroServiceUserLoginRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&userLoginRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.MicroServiceUserLoginResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	accessToken, err := mulc.microServiceServiceObj.GenerateMicroServiceUserAccessKey(userLoginRequest.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MicroServiceUserLoginResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, models.MicroServiceUserLoginResponse{
		Status:            1,
		StatusDescription: "Access Token generated",
		AccessToken:       accessToken,
	})
}
