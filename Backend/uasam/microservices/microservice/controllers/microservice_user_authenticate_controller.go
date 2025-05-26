package controllers

import (
	"net/http"
	"uasam/logger"
	"uasam/microservices/microservice/models"
	microServiceService "uasam/microservices/microservice/services"

	"github.com/gin-gonic/gin"
)

type MicroServiceUserAuthenticateController struct {
	microServiceServiceObj *microServiceService.MicroServiceService
	log                    *logger.Logger
}

func NewMicroServiceUserAuthenticateController(microServiceServiceObj *microServiceService.MicroServiceService, log *logger.Logger) *MicroServiceUserAuthenticateController {
	return &MicroServiceUserAuthenticateController{
		microServiceServiceObj: microServiceServiceObj,
		log:                    log,
	}
}

func (muac *MicroServiceUserAuthenticateController) MicroServiceUserAuthenticateHandler(c *gin.Context) {
	var header models.MicroServiceUserAuthenticateRequestHeaders
	if err := c.ShouldBindHeader(&header); err != nil {
		c.JSON(http.StatusBadRequest, models.MicroServiceUserAuthenticateResponse{
			Status:            -1,
			StatusDescription: "Missing or invalid required headers",
		})
		return
	}

	userID, err := muac.microServiceServiceObj.AuthenticateUserMicroService(header.UserToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MicroServiceUserAuthenticateResponse{
			Status:            -1,
			StatusDescription: "Not Authorized",
		})
		return
	}

	c.JSON(http.StatusOK, models.MicroServiceUserAuthenticateResponse{
		Status:            1,
		StatusDescription: "Token authorization success",
		UserID:            userID,
	})
}
