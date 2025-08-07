package controllers

import (
	"fmt"
	"governor/logger"
	"governor/user/algorithm/models"
	"governor/user/algorithm/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserAlgorithmController struct {
	log                  *logger.Logger
	userAlgorithmService *services.UserAlgorithmService
}

func NewUserAlgorithmController(log *logger.Logger, UserAlgorithmService *services.UserAlgorithmService) *UserAlgorithmController {
	return &UserAlgorithmController{
		log:                  log,
		userAlgorithmService: UserAlgorithmService,
	}
}

// CreateUserAlgorithmHandler uploads a new algorithm script for the authenticated user.
//
// @Summary      Create a user algorithm
// @Description  Uploads a script file and creates a new user algorithm associated with the authenticated user.
// @Tags         UserAlgorithm
// @Accept       multipart/form-data
// @Produce      json
// @Param        script_name formData string true "Script Name"
// @Param        script      formData file   true "Python script file"
// @Success      201 {object} models.CreateUserAlgorithmResponse "Algorithm created successfully"
// @Failure      400 {object} models.CreateUserAlgorithmResponse "Invalid input or file upload failed"
// @Failure      500 {object} models.CreateUserAlgorithmResponse "Internal server error"
// @Security     USER_TOKEN
// @Router       /v1/user/algorithm/python/upload/ [post]
func (uac *UserAlgorithmController) CreateUserAlgorithmHandler(c *gin.Context) {
	var createUserAlgorithmRequest models.CreateUserAlgorithmRequest
	userID, _ := c.Get("USER_ID")

	// Bind form fields (from multipart/form-data)
	if err := c.ShouldBind(&createUserAlgorithmRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.CreateUserAlgorithmResponse{
			Status:            -1,
			StatusDescription: "Invalid form values",
		})
		return
	}

	script, _, err := c.Request.FormFile("algorithm")
	if err != nil {
		go uac.log.Error("error in getting file from form", zap.String("execution level", "CreateUserAlgorithmHandler"), zap.String("Error", err.Error()))
		c.JSON(http.StatusBadRequest, models.CreateUserAlgorithmResponse{
			Status:            -1,
			StatusDescription: "File upload failed",
		})
		return
	}
	defer script.Close()

	userAlgorithm, err := uac.userAlgorithmService.CreateUserAlgorithmHandler(fmt.Sprint(userID), createUserAlgorithmRequest.AlgorithmName, script)
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

func (uac *UserAlgorithmController) EditUserAlgorithmHandler(c *gin.Context) {
	var editUserAlgorithmRequest models.EditUserAlgorithmRequest
	userID, _ := c.Get("USER_ID")

	// Bind form fields (from multipart/form-data)
	if err := c.ShouldBind(&editUserAlgorithmRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.EditUserAlgorithmResponse{
			Status:            -1,
			StatusDescription: "Invalid form values",
		})
		return
	}

	algorithmID, err := uuid.Parse(editUserAlgorithmRequest.AlgorithmID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.EditUserAlgorithmResponse{
			Status:            -1,
			StatusDescription: "Invalid algorithm ID",
		})
		return
	}

	script, _, err := c.Request.FormFile("algorithm")
	if err != nil {
		go uac.log.Error("error in getting file from form", zap.String("execution level", "EditUserAlgorithmHandler"), zap.String("Error", err.Error()))
		c.JSON(http.StatusBadRequest, models.EditUserAlgorithmResponse{
			Status:            -1,
			StatusDescription: "File upload failed",
		})
		return
	}
	defer script.Close()

	userAlgorithm, err := uac.userAlgorithmService.EditUserAlgorithmHandler(fmt.Sprint(userID), algorithmID, editUserAlgorithmRequest.AlgorithmName, script)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.EditUserAlgorithmResponse{
			Status:            0,
			StatusDescription: "Algorithm edit failed",
		})
		return
	}

	c.JSON(http.StatusOK, models.EditUserAlgorithmResponse{
		Status:            1,
		StatusDescription: "User Algorithm edited successfully",
		UserAlgorithm:     userAlgorithm,
	})
}
