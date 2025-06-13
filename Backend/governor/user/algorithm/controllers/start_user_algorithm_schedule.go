package controllers

import (
	"encoding/json"
	"fmt"
	"governor/user/algorithm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (uac *UserAlgorithmController) StartUserAlgorithmCronSchedule(c *gin.Context) {
	var startUserAlgorithmCronScheduleRequest models.StartUserAlgorithmCronScheduleRequest
	userID, _ := c.Get("USER_ID")

	if err := json.NewDecoder(c.Request.Body).Decode(&startUserAlgorithmCronScheduleRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.StartUserAlgorithmCronScheduleResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	err := uac.userAlgorithmService.StartUserAlgorithm(fmt.Sprint(userID), startUserAlgorithmCronScheduleRequest.AlgorithmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StartUserAlgorithmCronScheduleResponse{
			Status:            0,
			StatusDescription: "Algorithm schedule update failed",
		})
		return
	}

	c.JSON(http.StatusOK, models.StartUserAlgorithmCronScheduleResponse{
		Status:            1,
		StatusDescription: "User algorithm started successfully",
	})
}
