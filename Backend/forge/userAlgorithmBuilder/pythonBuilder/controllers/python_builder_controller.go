package controllers

import (
	"encoding/json"
	"forge/logger"
	"forge/userAlgorithmBuilder/pythonBuilder/models"
	"forge/userAlgorithmBuilder/pythonBuilder/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PythonBuilderController struct {
	log                  *logger.Logger
	pythonBuilderService *services.PythonBuilderService
}

func NewPythonBuilderController(log *logger.Logger, pythonBuilderService *services.PythonBuilderService) *PythonBuilderController {
	return &PythonBuilderController{
		log:                  log,
		pythonBuilderService: pythonBuilderService,
	}
}

// PythonBuilderHandler handles building and pushing a user-submitted Python algorithm.
// @Summary Build and Push Python Algorithm
// @Description Accepts a user ID, script ID, and script URL to build and push the Python algorithm as a Docker image.
// @Tags PythonBuilder
// @Accept json
// @Produce json
// @Param PythonBuilderRequest body models.PythonBuilderRequest true "Python Builder Request"
// @Success 201 {object} models.PythonBuilderResponse "User Algorithm Built and Pushed successfully"
// @Failure 400 {object} models.PythonBuilderResponse "Invalid request payload"
// @Failure 500 {object} models.PythonBuilderResponse "Internal Server Error"
// @Router /v1/algorithm/python/build/ [post]
func (pbc *PythonBuilderController) PythonBuilderHandler(c *gin.Context) {
	var pythonBuilderRequest models.PythonBuilderRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&pythonBuilderRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.PythonBuilderResponse{
			Status:            -1,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	err := pbc.pythonBuilderService.BuildAlgorithmHandler(pythonBuilderRequest.UserID, pythonBuilderRequest.ScriptID, pythonBuilderRequest.ScriptURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PythonBuilderResponse{
			Status:            0,
			StatusDescription: "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusCreated, models.PythonBuilderResponse{
		Status:            1,
		StatusDescription: "User Algorithm Built and Pushed successfully",
	})
}
