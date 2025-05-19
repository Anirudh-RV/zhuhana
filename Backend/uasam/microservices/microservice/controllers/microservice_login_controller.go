package controllers

import (
	"net/http"
	"uasam/logger"
	"uasam/microservices/microservice/models"
	microServiceService "uasam/microservices/microservice/services"

	"github.com/gin-gonic/gin"
)

type MicroServiceLoginController struct {
	microServiceServiceObj *microServiceService.MicroServiceService
	log                    *logger.Logger
}

func NewMicroServiceLoginController(microServiceServiceObj *microServiceService.MicroServiceService, log *logger.Logger) *MicroServiceLoginController {
	return &MicroServiceLoginController{
		microServiceServiceObj: microServiceServiceObj,
		log:                    log,
	}
}

// MicroServiceLoginHandler godoc
// @Summary Authenticate microservice and generate access token
// @Description Validates the incoming request headers from a calling microservice and generates an access token if valid
// @Tags Microservice
// @Accept  json
// @Produce  json
// @Param Caller-Service header string true "Name of the calling microservice"
// @Param API-Key header string true "API key for the calling microservice"
// @Success 200 {object} models.MicroServiceLoginResponse "Access token generated successfully"
// @Failure 400 {object} models.MicroServiceLoginResponse "Missing or invalid required headers"
// @Failure 500 {object} models.MicroServiceLoginResponse "Internal server error while generating token"
// @Router /v1/microservice/login [post]
func (mlc *MicroServiceLoginController) MicroServiceLoginHandler(c *gin.Context) {
	var header models.MicroServiceLoginRequestHeaders
	if err := c.ShouldBindHeader(&header); err != nil {
		c.JSON(http.StatusBadRequest, models.MicroServiceLoginResponse{
			Status:            -1,
			StatusDescription: "Missing or invalid required headers",
		})
		return
	}

	accessToken, err := mlc.microServiceServiceObj.GenerateMicroServiceAccessKey(header.CallerService, header.APIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MicroServiceLoginResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, models.MicroServiceLoginResponse{
		Status:            1,
		StatusDescription: "Access Token generated",
		AccessToken:       accessToken,
	})
}
