package controllers

import (
	"encoding/json"
	"fmt"
	"governor/user/algorithm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StopUserAlgorithm godoc
//
// @Summary      Stops the user algorithm using a cron schedule
// @Description  Initiates the scheduled execution of a user algorithm based on the provided algorithm ID.
// @Tags         UserAlgorithm
// @Accept       json
// @Produce      json
// @Param        request   body      models.StopUserAlgorithmRequest  true  "Stop Algorithm Cron Request"
// @Success      200       {object}  models.StopUserAlgorithmResponse
// @Failure      400       {object}  models.StopUserAlgorithmResponse
// @Failure      500       {object}  models.StopUserAlgorithmResponse
// @Security     USER_TOKEN
// @Router       /user/algorithm/start/ [post]
func (uac *UserAlgorithmController) StopUserAlgorithm(c *gin.Context) {
	var stopUserAlgorithmRequest models.StopUserAlgorithmRequest
	userID, _ := c.Get("USER_ID")

	if err := json.NewDecoder(c.Request.Body).Decode(&stopUserAlgorithmRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.StopUserAlgorithmResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	err := uac.userAlgorithmService.StopUserAlgorithm(fmt.Sprint(userID), stopUserAlgorithmRequest.AlgorithmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StopUserAlgorithmResponse{
			Status:            0,
			StatusDescription: "Algorithm schedule update failed",
		})
		return
	}

	c.JSON(http.StatusOK, models.StopUserAlgorithmResponse{
		Status:            1,
		StatusDescription: "User algorithm stopped successfully",
	})
}
