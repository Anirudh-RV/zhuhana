package controllers

import (
	"encoding/json"
	"net/http"
	"uasam/logger"
	"uasam/users/user/models"
	userService "uasam/users/user/services"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	userService *userService.UserService
	log         *logger.Logger
}

func NewLoginController(userService *userService.UserService, log *logger.Logger) *LoginController {
	return &LoginController{
		userService: userService,
		log:         log,
	}
}

// LoginVerifyPasswordHandler godoc
// @Summary Verify user password during login
// @Description Checks if user exists and verifies the submitted password. Sends an OTP if the password verification is successful
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body models.LoginVerifyPasswordRequest true "Login Verify Password Request"
// @Success 200 {object} models.LoginVerifyPasswordResponse "Password verified successfully, proceed to OTP verification"
// @Failure 400 {object} models.LoginVerifyPasswordResponse "Invalid request payload"
// @Failure 401 {object} models.LoginVerifyPasswordResponse "Login error due to invalid credentials or user not existing"
// @Failure 500 {object} models.LoginVerifyPasswordResponse "Internal server error"
// @Router /v1/user/login/verify-password/ [post]
func (lgc *LoginController) LoginVerifyPasswordHandler(c *gin.Context) {
	var loginVerifyPasswordRequest models.LoginVerifyPasswordRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&loginVerifyPasswordRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.LoginVerifyPasswordResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	userExists, err := lgc.userService.IfUserExists(loginVerifyPasswordRequest.EmailID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.LoginVerifyPasswordResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}
	if !userExists {
		c.JSON(http.StatusUnauthorized, models.LoginVerifyPasswordResponse{
			Status:            -1,
			StatusDescription: "Login Error",
		})
		return
	}

	err = lgc.userService.LoginVerifyPasswordHandler(&loginVerifyPasswordRequest, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.LoginVerifyPasswordResponse{
			Status:            0,
			StatusDescription: "Login Error",
		})
		return
	}

	c.JSON(http.StatusOK, &models.LoginVerifyPasswordResponse{
		Status:            1,
		StatusDescription: "Verify OTP",
	})
}

// LoginVerifyOTPHandler godoc
// @Summary Verify OTP during login
// @Description Verifies the OTP provided by the user and returns an access token upon successful authentication
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body models.LoginVerifyOTPRequest true "Login Verify OTP Request"
// @Success 200 {object} models.LoginVerifyOTPResponse "Login authenticated successfully, access token issued"
// @Failure 400 {object} models.LoginVerifyOTPResponse "Invalid request payload or login error"
// @Failure 401 {object} models.LoginVerifyOTPResponse "User does not exist or login error"
// @Failure 500 {object} models.LoginVerifyOTPResponse "Internal server error"
// @Router /v1/user/login/verify-otp/ [post]
func (lgc *LoginController) LoginVerifyOTPHandler(c *gin.Context) {
	var loginVerifyOTPRequest models.LoginVerifyOTPRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&loginVerifyOTPRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.LoginVerifyOTPResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	userExists, err := lgc.userService.IfUserExists(loginVerifyOTPRequest.EmailID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.LoginVerifyOTPResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}
	if !userExists {
		c.JSON(http.StatusUnauthorized, models.LoginVerifyOTPResponse{
			Status:            -1,
			StatusDescription: "Login Error",
		})
		return
	}

	userResponseObject, generatedUserAccessToken, err := lgc.userService.LoginVerifyOTPHandler(&loginVerifyOTPRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.LoginVerifyOTPResponse{
			Status:            -1,
			StatusDescription: "Login Error",
		})
		return
	}

	c.JSON(http.StatusOK, &models.LoginVerifyOTPResponse{
		Status:            1,
		StatusDescription: "Login Authenticated",
		User:              *userResponseObject,
		AccessToken:       generatedUserAccessToken,
	})
}
