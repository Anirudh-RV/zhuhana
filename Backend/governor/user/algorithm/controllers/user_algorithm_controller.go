package controllers

import (
	"fmt"
	"governor/logger"
	"governor/user/algorithm/models"
	"governor/user/algorithm/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAlgorithmController struct {
	log                  *logger.Logger
	UserAlgorithmService *services.UserAlgorithmService
}

func NewUserAlgorithmController(log *logger.Logger, UserAlgorithmService *services.UserAlgorithmService) *UserAlgorithmController {
	return &UserAlgorithmController{
		log:                  log,
		UserAlgorithmService: UserAlgorithmService,
	}
}

func (uac *UserAlgorithmController) CreateUserAlgorithmHandler(c *gin.Context) {
	var createUserAlgorithmRequest models.CreateUserAlgorithmRequest

	// Bind form fields (from multipart/form-data)
	if err := c.ShouldBind(&createUserAlgorithmRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.CreateUserAlgorithmResponse{
			Status:            -1,
			StatusDescription: "Invalid form values",
		})
		return
	}

	script, _, err := c.Request.FormFile("script")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.CreateUserAlgorithmResponse{
			Status:            -1,
			StatusDescription: "File upload failed",
		})
		return
	}
	defer script.Close()

	userID, _ := c.Get("USER_ID")
	userAlgorithm, err := uac.UserAlgorithmService.UserAlgorithmHandler(fmt.Sprint(userID), createUserAlgorithmRequest.ScriptName, createUserAlgorithmRequest.CronSchedule, script)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.CreateUserAlgorithmResponse{
			Status:            0,
			StatusDescription: "Algorithm creation failed",
		})
		return
	}

	c.JSON(http.StatusCreated, models.CreateUserAlgorithmResponse{
		Status:            1,
		StatusDescription: "User Algorithm uploaded successfully",
		UserAlgorithm:     userAlgorithm,
	})
}
