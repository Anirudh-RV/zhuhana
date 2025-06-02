package controllers

import (
	"net/http"
	"uasam/logger"
	"uasam/users/user/models"
	"uasam/users/user/services"

	"github.com/gin-gonic/gin"
)

type UserAuthenticateController struct {
	userService *services.UserService
	log         *logger.Logger
}

func NewUserAuthenticateController(log *logger.Logger, userService *services.UserService) *UserAuthenticateController {
	return &UserAuthenticateController{
		userService: userService,
		log:         log,
	}
}

func (uac *UserAuthenticateController) UserAuthenticateHandler(c *gin.Context) {
	var header models.UserAuthenticateRequestHeaders
	if err := c.ShouldBindHeader(&header); err != nil {
		c.JSON(http.StatusBadRequest, models.UserAuthenticateResponse{
			Status:            -1,
			StatusDescription: "Missing or invalid required headers",
		})
		return
	}

	userObjPtr, err := uac.userService.AuthenticateUser(header.UserToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.UserAuthenticateResponse{
			Status:            -1,
			StatusDescription: "Not Authorized",
		})
		return
	}

	c.JSON(http.StatusOK, models.UserAuthenticateResponse{
		Status:            1,
		StatusDescription: "Token authentication success",
		User:              userObjPtr,
	})
}
