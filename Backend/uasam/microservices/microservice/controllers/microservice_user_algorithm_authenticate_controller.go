package controllers

import (
	"net/http"
	"uasam/logger"
	"uasam/microservices/microservice/models"
	microServiceService "uasam/microservices/microservice/services"

	"github.com/gin-gonic/gin"
)

type MicroServiceUserAlgorithmAuthenticateController struct {
	microServiceServiceObj *microServiceService.MicroServiceService
	log                    *logger.Logger
}

func NewMicroServiceUserAlgorithmAuthenticateController(microServiceServiceObj *microServiceService.MicroServiceService, log *logger.Logger) *MicroServiceUserAlgorithmAuthenticateController {
	return &MicroServiceUserAlgorithmAuthenticateController{
		microServiceServiceObj: microServiceServiceObj,
		log:                    log,
	}
}

// MicroServiceUserAlgorithmAuthenticateHandler godoc
// @Summary Authorize a microservice user
// @Description Validates the access token of a microservice user algorithm and returns the user algorithm ID
// @Tags Microservice Authentication
// @Accept json
// @Produce json
// @Param authHeaders header models.MicroServiceUserAlgorithmAuthenticateRequestHeaders true "Access Token of the microservice user algorithm"
// @Success 200 {object} models.MicroServiceUserAlgorithmAuthenticateResponse
// @Failure 400 {object} models.MicroServiceUserAuthenticateResponse "Missing or invalid required headers"
// @Failure 401 {object} models.MicroServiceUserAuthenticateResponse "Not Authorized"
// @Router /v1/microservice/user/authenticate/ [get]
func (muac *MicroServiceUserAlgorithmAuthenticateController) MicroServiceUserAlgorithmAuthenticateHandler(c *gin.Context) {
	var header models.MicroServiceUserAlgorithmAuthenticateRequestHeaders
	if err := c.ShouldBindHeader(&header); err != nil {
		c.JSON(http.StatusBadRequest, models.MicroServiceUserAlgorithmAuthenticateResponse{
			Status:            -1,
			StatusDescription: "Missing or invalid required headers",
		})
		return
	}

	userAlgorithmID, err := muac.microServiceServiceObj.AuthenticateMicroServiceUserAlgorithm(header.UserAlgorithmToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MicroServiceUserAlgorithmAuthenticateResponse{
			Status:            -1,
			StatusDescription: "Not Authorized",
		})
		return
	}

	c.JSON(http.StatusOK, models.MicroServiceUserAlgorithmAuthenticateResponse{
		Status:            1,
		StatusDescription: "Token authentication success",
		UserAlgorithmID:   userAlgorithmID,
	})
}
