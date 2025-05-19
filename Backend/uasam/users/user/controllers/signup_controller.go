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
		c.JSON(http.StatusBadRequest, models.SignUpInitResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
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

	err = snc.userService.SignUpInitHandler(&signUpInitRequest, c.Request.UserAgent(), c.ClientIP())
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

// SignUpVerifyOTPHandler verifies the OTP sent to the user and creates the user account.
// @Summary      Verify OTP and Sign Up User
// @Description  Verifies the OTP and creates the user account if valid.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request  body      models.SignUpVerifyOTPRequest  true  "Verify OTP Request"
// @Success      201      {object}  models.SignUpVerifyOTPResponse
// @Failure      400      {object}  models.SignUpVerifyOTPResponse  "Invalid request payload"
// @Failure      401      {object}  models.SignUpVerifyOTPResponse  "Wrong OTP"
// @Failure      500      {object}  models.SignUpVerifyOTPResponse       "Server Error"
// @Router       /v1/user/sign-up/verify-otp/ [post]
func (snc *SignUpController) SignUpVerifyOTPHandler(c *gin.Context) {
	var signUpVerifyOTPRequest models.SignUpVerifyOTPRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&signUpVerifyOTPRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.SignUpVerifyOTPResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	userExists, err := snc.userService.IfUserExists(signUpVerifyOTPRequest.EmailID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.SignUpVerifyOTPResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}
	if userExists {
		c.JSON(http.StatusBadRequest, models.SignUpVerifyOTPResponse{
			Status:            -1,
			StatusDescription: "User with Email Already Exists",
		})
		return
	}

	userResponseObject, generatedUserAccessToken, status, err := snc.userService.SignUpVerifyOTPHandler(&signUpVerifyOTPRequest)
	if status == -1 {
		c.JSON(http.StatusUnauthorized, models.SignUpVerifyOTPResponse{
			Status:            -2,
			StatusDescription: "Wrong OTP",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.SignUpVerifyOTPResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}

	c.JSON(http.StatusCreated, &models.SignUpVerifyOTPResponse{
		Status:            1,
		StatusDescription: "user created",
		User:              *userResponseObject,
		AccessToken:       generatedUserAccessToken,
	})
}
