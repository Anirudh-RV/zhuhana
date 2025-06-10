package controllers

import (
	"encoding/json"
	"fmt"
	"governor/user/algorithm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CancelUserAlgorithmCronSchedule cancels the start and end cron schedules for a user algorithm.
//
// @Summary Cancel cron schedule for a user's algorithm
// @Description Cancels the start and end cron jobs of a specific user algorithm, if the algorithm belongs to the authenticated user.
// @Tags UserAlgorithm
// @Accept json
// @Produce json
// @Param request body models.CancelUserAlgorithmCronScheduleRequest true "Algorithm cancellation request"
// @Success 200 {object} models.CancelUserAlgorithmCronScheduleResponse "Successfully canceled algorithm schedule"
// @Failure 400 {object} models.CancelUserAlgorithmCronScheduleResponse "Invalid request payload"
// @Failure 500 {object} models.CancelUserAlgorithmCronScheduleResponse "Internal server error during schedule cancellation"
// @Router /v1/user/algorithm/schedule/cancel/ [post]
// @Security USER_TOKEN
func (uac *UserAlgorithmController) CancelUserAlgorithmCronSchedule(c *gin.Context) {
	var cancelUserAlgorithmCronScheduleRequest models.CancelUserAlgorithmCronScheduleRequest
	userID, _ := c.Get("USER_ID")

	if err := json.NewDecoder(c.Request.Body).Decode(&cancelUserAlgorithmCronScheduleRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.CancelUserAlgorithmCronScheduleResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	err := uac.userAlgorithmService.CancelAlgorithmSchedule(fmt.Sprint(userID), cancelUserAlgorithmCronScheduleRequest.AlgorithmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.CancelUserAlgorithmCronScheduleResponse{
			Status:            0,
			StatusDescription: "Algorithm schedule update failed",
		})
		return
	}

	c.JSON(http.StatusOK, models.CancelUserAlgorithmCronScheduleResponse{
		Status:            1,
		StatusDescription: "User algorithm schedule updated successfully",
	})
}
