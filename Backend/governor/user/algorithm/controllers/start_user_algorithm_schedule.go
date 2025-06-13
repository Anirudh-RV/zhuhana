package controllers

import (
	"encoding/json"
	"fmt"
	"governor/user/algorithm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StartUserAlgorithm godoc
//
// @Summary      Starts the user algorithm using a cron schedule
// @Description  Initiates the scheduled execution of a user algorithm based on the provided algorithm ID.
// @Tags         UserAlgorithm
// @Accept       json
// @Produce      json
// @Param        request   body      models.StartUserAlgorithmRequest  true  "Start Algorithm Cron Request"
// @Success      200       {object}  models.StartUserAlgorithmResponse
// @Failure      400       {object}  models.StartUserAlgorithmResponse
// @Failure      500       {object}  models.StartUserAlgorithmResponse
// @Security     USER_TOKEN
// @Router       /user/algorithm/start/ [post]
func (uac *UserAlgorithmController) StartUserAlgorithm(c *gin.Context) {
	var startUserAlgorithmRequest models.StartUserAlgorithmRequest
	userID, _ := c.Get("USER_ID")

	if err := json.NewDecoder(c.Request.Body).Decode(&startUserAlgorithmRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.StartUserAlgorithmResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	err := uac.userAlgorithmService.StartUserAlgorithm(fmt.Sprint(userID), startUserAlgorithmRequest.AlgorithmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StartUserAlgorithmResponse{
			Status:            0,
			StatusDescription: "Algorithm schedule update failed",
		})
		return
	}

	c.JSON(http.StatusOK, models.StartUserAlgorithmResponse{
		Status:            1,
		StatusDescription: "User algorithm started successfully",
	})
}
