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

func (muc *MicroServiceUserLoginController) MicroServiceUserLoginHandler(c *gin.Context) {
	var userLoginRequest models.MicroServiceUserLoginRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&userLoginRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.MicroServiceUserLoginResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	accessToken, err := muc.microServiceServiceObj.GenerateMicroServiceUserAccessKey(userLoginRequest.UserID)
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
