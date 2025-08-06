package controllers

import (
	"fmt"
	"governor/user/algorithm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (uac *UserAlgorithmController) GetUserAlgorithmRunsHandler(c *gin.Context) {
	userID, _ := c.Get("USER_ID")
	algorithmID := c.Query("algorithm_id")

	if algorithmID == "" {
		c.JSON(http.StatusBadRequest, models.GetAllUserAlgorithmRunsResponse{
			Status:            0,
			StatusDescription: "Missing query param: algorithm_id",
		})
		return
	}

	userAlgorithmRuns, err := uac.userAlgorithmService.GetUserAlgorithmRuns(fmt.Sprint(userID), fmt.Sprint(algorithmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.GetAllUserAlgorithmRunsResponse{
			Status:            0,
			StatusDescription: "Getting user algorithm runs failed",
		})
		return
	}

	c.JSON(http.StatusOK, models.GetAllUserAlgorithmRunsResponse{
		Status:            1,
		StatusDescription: "Fetched user algorithm runs successfully",
		UserAlgorithmRuns: userAlgorithmRuns,
	})
}
