package controllers

import (
	"encoding/json"
	"fmt"
	"governor/user/algorithm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateUserAlgorithmCronSchedule updates the cron schedule for a specific user algorithm.
//
// @Summary      Update user algorithm cron schedule
// @Description  Updates the cron schedule for an existing user algorithm belonging to the authenticated user.
// @Tags         UserAlgorithm
// @Accept       json
// @Produce      json
// @Param        body body models.UpdateUserAlgorithmCronScheduleRequest true "Cron schedule update payload"
// @Success      200 {object} models.UpdateUserAlgorithmCronScheduleResponse "Cron schedule updated successfully"
// @Failure      400 {object} models.UpdateUserAlgorithmCronScheduleResponse "Invalid request payload"
// @Failure      500 {object} models.UpdateUserAlgorithmCronScheduleResponse "Internal server error"
// @Security     USER_TOKEN
// @Router       /v1/user/algorithm/schedule/ [post]
func (uac *UserAlgorithmController) UpdateUserAlgorithmCronSchedule(c *gin.Context) {
	var updateUserAlgorithmCronScheduleRequest models.UpdateUserAlgorithmCronScheduleRequest
	userID, _ := c.Get("USER_ID")

	if err := json.NewDecoder(c.Request.Body).Decode(&updateUserAlgorithmCronScheduleRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.UpdateUserAlgorithmCronScheduleResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	err := uac.userAlgorithmService.UpdateAlgorithmSchedule(fmt.Sprint(userID), updateUserAlgorithmCronScheduleRequest.AlgorithmID, updateUserAlgorithmCronScheduleRequest.StartCronSchedule, updateUserAlgorithmCronScheduleRequest.EndCronSchedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.UpdateUserAlgorithmCronScheduleResponse{
			Status:            0,
			StatusDescription: "Algorithm schedule update failed",
		})
		return
	}

	c.JSON(http.StatusOK, models.UpdateUserAlgorithmCronScheduleResponse{
		Status:            1,
		StatusDescription: "User algorithm schedule updated successfully",
	})
}
