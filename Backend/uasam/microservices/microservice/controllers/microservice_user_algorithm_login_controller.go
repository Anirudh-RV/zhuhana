package controllers

import (
	"encoding/json"
	"net/http"
	"uasam/logger"
	"uasam/microservices/microservice/models"
	microServiceService "uasam/microservices/microservice/services"

	"github.com/gin-gonic/gin"
)

type MicroServiceUserAlgorithmLoginController struct {
	microServiceServiceObj *microServiceService.MicroServiceService
	log                    *logger.Logger
}

func NewMicroServiceUserAlgorithmLoginController(microServiceServiceObj *microServiceService.MicroServiceService, log *logger.Logger) *MicroServiceUserAlgorithmLoginController {
	return &MicroServiceUserAlgorithmLoginController{
		microServiceServiceObj: microServiceServiceObj,
		log:                    log,
	}
}

// MicroServiceUserAlgorithmLoginHandler godoc
// @Summary Authenticate a microservice user algorithm
// @Description Generates an access token for a valid microservice user algorithm
// @Tags Microservice Authentication
// @Accept json
// @Produce json
// @Param user body models.MicroServiceUserAlgorithmLoginRequest true "Microservice User Algorithm Login Request"
// @Success 200 {object} models.MicroServiceUserAlgorithmLoginResponse
// @Failure 400 {object} models.MicroServiceUserAlgorithmLoginResponse "Invalid request payload"
// @Failure 500 {object} models.MicroServiceUserAlgorithmLoginResponse "Server error"
// @Router /v1/microservice/user/login/ [post]
func (mulc *MicroServiceUserAlgorithmLoginController) MicroServiceUserAlgorithmLoginHandler(c *gin.Context) {
	var userAlgorithmLoginRequest models.MicroServiceUserAlgorithmLoginRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&userAlgorithmLoginRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.MicroServiceUserAlgorithmLoginResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	accessToken, err := mulc.microServiceServiceObj.GenerateMicroServiceUserAlgorithmAccessKey(userAlgorithmLoginRequest.UserAlgorithmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MicroServiceUserAlgorithmLoginResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, models.MicroServiceUserAlgorithmLoginResponse{
		Status:            1,
		StatusDescription: "Access Token generated",
		AccessToken:       accessToken,
	})
}
