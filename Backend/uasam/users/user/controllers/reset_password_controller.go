package controllers

import (
	"encoding/json"
	"net/http"
	"uasam/logger"
	"uasam/users/user/models"
	userService "uasam/users/user/services"

	"github.com/gin-gonic/gin"
)

type ResetPasswordController struct {
	userService *userService.UserService
	log         *logger.Logger
}

func NewResetPasswordController(userService *userService.UserService, log *logger.Logger) *ResetPasswordController {
	return &ResetPasswordController{
		userService: userService,
		log:         log,
	}
}

// ResetPasswordInitHandler godoc
// @Summary Initiate Password Reset
// @Description Starts the password reset process by verifying if the user exists and sending a reset link to the email
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body models.ResetPasswordInitRequest true "Reset Password Init Request"
// @Success 200 {object} models.ResetPasswordInitResponse "Password reset link sent successfully"
// @Failure 400 {object} models.ResetPasswordInitResponse "Invalid request payload"
// @Failure 401 {object} models.ResetPasswordInitResponse "User does not exist or reset error"
// @Failure 500 {object} models.ResetPasswordInitResponse "Server error"
// @Router /v1/user/reset-password/init/ [post]
func (rpc *ResetPasswordController) ResetPasswordInitHandler(c *gin.Context) {
	var resetPasswordInitRequest models.ResetPasswordInitRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&resetPasswordInitRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.ResetPasswordInitResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	userExists, err := rpc.userService.IfUserExists(resetPasswordInitRequest.EmailID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResetPasswordInitResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}
	if !userExists {
		c.JSON(http.StatusUnauthorized, models.ResetPasswordInitResponse{
			Status:            -1,
			StatusDescription: "Reset password Error",
		})
		return
	}

	err = rpc.userService.ResetPasswordInitHandler(&resetPasswordInitRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ResetPasswordInitResponse{
			Status:            0,
			StatusDescription: "Reset password Error",
		})
		return
	}

	c.JSON(http.StatusOK, &models.ResetPasswordInitResponse{
		Status:            1,
		StatusDescription: "Check email for reset password link",
	})
}

// ResetPasswordHandler godoc
// @Summary Complete Password Reset
// @Description Completes the password reset by setting a new password for a valid user after token authentication
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body models.ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} models.ResetPasswordResponse "Password reset successful"
// @Failure 400 {object} models.ResetPasswordResponse "Invalid request payload"
// @Failure 401 {object} models.ResetPasswordResponse "User does not exist or reset error"
// @Failure 500 {object} models.ResetPasswordResponse "Server error"
// @Router /v1/user/reset-password/reset/ [post]
func (rpc *ResetPasswordController) ResetPasswordHandler(c *gin.Context) {
	var resetPasswordRequest models.ResetPasswordRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&resetPasswordRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.ResetPasswordResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	userExists, err := rpc.userService.IfUserExists(resetPasswordRequest.EmailID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResetPasswordResponse{
			Status:            0,
			StatusDescription: "Server Error",
		})
		return
	}
	if !userExists {
		c.JSON(http.StatusUnauthorized, models.ResetPasswordResponse{
			Status:            -1,
			StatusDescription: "Reset password Error",
		})
		return
	}

	err = rpc.userService.ResetPasswordHandler(&resetPasswordRequest)

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ResetPasswordResponse{
			Status:            -1,
			StatusDescription: "Reset password Error",
		})
		return
	}

	c.JSON(http.StatusOK, models.ResetPasswordResponse{
		Status:            1,
		StatusDescription: "Reset password successful",
	})
}
