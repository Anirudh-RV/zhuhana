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

// MicroServiceUserAuthenticateHandler godoc
// @Summary Authorize a microservice user
// @Description Validates the access token of a microservice user and returns the user ID
// @Tags Microservice Authentication
// @Accept json
// @Produce json
// @Param authHeaders header models.MicroServiceUserAuthenticateRequestHeaders true "Access Token of the microservice user"
// @Success 200 {object} models.MicroServiceUserAuthenticateResponse
// @Failure 400 {object} models.MicroServiceUserAuthenticateResponse "Missing or invalid required headers"
// @Failure 401 {object} models.MicroServiceUserAuthenticateResponse "Not Authorized"
// @Router /v1/microservice/user/authenticate/ [get]
func (muac *MicroServiceUserAuthenticateController) MicroServiceUserAuthenticateHandler(c *gin.Context) {
	var header models.MicroServiceUserAuthenticateRequestHeaders
	if err := c.ShouldBindHeader(&header); err != nil {
		c.JSON(http.StatusBadRequest, models.MicroServiceUserAuthenticateResponse{
			Status:            -1,
			StatusDescription: "Missing or invalid required headers",
		})
		return
	}

	userID, err := muac.microServiceServiceObj.AuthenticateMicroServiceUserService(header.UserScriptToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MicroServiceUserAuthenticateResponse{
			Status:            -1,
			StatusDescription: "Not Authorized",
		})
		return
	}

	c.JSON(http.StatusOK, models.MicroServiceUserAuthenticateResponse{
		Status:            1,
		StatusDescription: "Token authentication success",
		UserID:            userID,
	})
}
