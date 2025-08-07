package controllers

import (
	"fmt"
	"governor/user/algorithm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserAlgorithms retrieves all user algorithms for the authenticated user.
//
// @Summary      Get all user algorithms
// @Description  Retrieves the list of all algorithms associated with the authenticated user.
// @Tags         UserAlgorithm
// @Produce      json
// @Success      200 {object} models.GetAllUserAlgorithmsResponse "User algorithms retrieved successfully"
// @Failure      500 {object} models.GetAllUserAlgorithmsResponse "Internal server error"
// @Security     USER_TOKEN
// @Router       /v1/user/algorithm/ [get]
func (uac *UserAlgorithmController) GetUserAlgorithms(c *gin.Context) {
	userID, _ := c.Get("USER_ID")

	userAlgorithms, err := uac.userAlgorithmService.GetAllUserAlgorithms(fmt.Sprint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.GetAllUserAlgorithmsResponse{
			Status:            0,
			StatusDescription: "Getting all user algorithms failed",
		})
		return
	}

	c.JSON(http.StatusOK, models.GetAllUserAlgorithmsResponse{
		Status:            1,
		StatusDescription: "Fetched all user algorithms successfully",
		UserAlgorithms:    userAlgorithms,
	})
}

// GetUserAlgorithmByID godoc
//
// @Summary      Get user algorithm details
// @Description  Retrieves a specific user algorithm by its ID.
// @Tags         UserAlgorithm
// @Accept       json
// @Produce      json
// @Param        algorithm_id  query     string  true  "Algorithm ID"
// @Success      200           {object}  models.GetUserAlgorithmResponse  "Algorithm fetched successfully"
// @Failure      400           {object}  models.GetUserAlgorithmResponse  "Missing or invalid algorithm_id"
// @Failure      500           {object}  models.GetUserAlgorithmResponse  "Internal server error"
// @Security     ApiKeyAuth
// @Router       /v1/user/algorithm/info [get]
func (uac *UserAlgorithmController) GetUserAlgorithmByID(c *gin.Context) {
	userID, _ := c.Get("USER_ID")
	algorithmID := c.Query("algorithm_id")

	if algorithmID == "" {
		c.JSON(http.StatusBadRequest, models.GetUserAlgorithmResponse{
			Status:            0,
			StatusDescription: "Missing query param: algorithm_id",
		})
		return
	}

	userAlgorithm, err := uac.userAlgorithmService.GetUserAlgorithm(fmt.Sprint(userID), fmt.Sprint(algorithmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.GetUserAlgorithmResponse{
			Status:            0,
			StatusDescription: "Getting user algorithm failed",
		})
		return
	}

	c.JSON(http.StatusOK, models.GetUserAlgorithmResponse{
		Status:            1,
		StatusDescription: "Fetched user algorithm successfully",
		UserAlgorithm:     userAlgorithm,
	})
}
