package controllers

import (
	"encoding/json"
	"net/http"
	"uasam/logger"
	"uasam/users/user/models"
	userService "uasam/users/user/services"

	"github.com/gin-gonic/gin"
)

type SignUpController struct {
	userService *userService.UserService
	log         *logger.Logger
}

func NewSignUpController(userService *userService.UserService, log *logger.Logger) *SignUpController {
	return &SignUpController{
		userService: userService,
		log:         log,
	}
}

// SignUpInitHandler godoc
// @Summary Initialize Sign Up
// @Description Initiates user sign-up by checking if the user exists and sending OTP for verification
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body models.SignUpInitRequest true "Sign Up Init Request"
// @Success 200 {object} models.SignUpInitResponse "OTP verification initiated"
// @Failure 400 {object} models.SignUpInitResponse "User already exists or invalid payload"
// @Failure 500 {object} models.SignUpInitResponse "Server error"
// @Router /v1/user/sign-up/init/ [post]
func (snc *SignUpController) SignUpInitHandler(c *gin.Context) {
	var signUpInitRequest models.SignUpInitRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&signUpInitRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":             0,
			"status_description": "Invalid request payload",
		})
		return
	}

	status, err := snc.userService.IfUserExists(signUpInitRequest.EmailID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.SignUpInitResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}

	if status {
		c.JSON(http.StatusBadRequest, models.SignUpInitResponse{
			Status:            -1,
			StatusDescription: "User with Email Already Exists",
		})
		return
	}

	err = snc.userService.SignUpInitHandler(&signUpInitRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.SignUpInitResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, models.SignUpInitResponse{
		Status:            1,
		StatusDescription: "verify OTP",
	})
}

func (snc *SignUpController) VerifyOTPHandler(c *gin.Context) {

}
