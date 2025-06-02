package controllers

import (
	"net/http"
	"uasam/logger"
	"uasam/microservices/microservice/models"
	microServiceService "uasam/microservices/microservice/services"

	"github.com/gin-gonic/gin"
)

type MicroServiceAuthenticateController struct {
	microServiceServiceObj *microServiceService.MicroServiceService
	log                    *logger.Logger
}

func NewMicroServiceAuthenticateController(microServiceServiceObj *microServiceService.MicroServiceService, log *logger.Logger) *MicroServiceAuthenticateController {
	return &MicroServiceAuthenticateController{
		microServiceServiceObj: microServiceServiceObj,
		log:                    log,
	}
}

// MicroServiceAuthenticateHandler godoc
// @Summary Authenticate microservice using JWT token
// @Description Authenticates the calling microservice using a JWT token passed in the headers. Returns the caller and callee service names on success.
// @Tags Microservice
// @Accept  json
// @Produce  json
// @Param OriginService header string true "Name of the calling microservice"
// @Param AuthToken header string true "JWT token for authentication"
// @Success 200 {object} models.MicroServiceAuthenticateResponse "Token authentication success"
// @Failure 400 {object} models.MicroServiceAuthenticateResponse "Missing or invalid required headers"
// @Failure 401 {object} models.MicroServiceAuthenticateResponse "Not Authorized"
// @Router /v1/microservice/authenticate/ [post]
func (mac *MicroServiceAuthenticateController) MicroServiceAuthenticateHandler(c *gin.Context) {
	var header models.MicroServiceAuthenticateRequestHeaders
	if err := c.ShouldBindHeader(&header); err != nil {
		c.JSON(http.StatusBadRequest, models.MicroServiceAuthenticateResponse{
			Status:            -1,
			StatusDescription: "Missing or invalid required headers",
		})
		return
	}

	callerService, calleeService, err := mac.microServiceServiceObj.AuthenticateMicroService(header.OriginService, header.AuthToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MicroServiceAuthenticateResponse{
			Status:            -1,
			StatusDescription: "Not Authorized",
		})
		return
	}

	c.JSON(http.StatusOK, models.MicroServiceAuthenticateResponse{
		Status:            1,
		StatusDescription: "Token authentication success",
		CallerService:     callerService,
		CalleeService:     calleeService,
	})
}
