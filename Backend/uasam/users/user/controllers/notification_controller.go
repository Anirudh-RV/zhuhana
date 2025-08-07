package controllers

import (
	"encoding/json"
	"net/http"
	"uasam/logger"
	"uasam/users/user/models"
	"uasam/users/user/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationController struct {
	notificationService *services.NotificationService
	log                 *logger.Logger
}

func NewNotificationController(log *logger.Logger, notificationService *services.NotificationService) *NotificationController {
	return &NotificationController{
		notificationService: notificationService,
		log:                 log,
	}
}

// GetNotificationsHandler godoc
//
// @Summary      Get user notifications
// @Description  Retrieves all notifications for the authenticated user.
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.GetNotificationsResponseSwagger  "Notifications fetched successfully"
// @Failure      400  {object}  map[string]string                 "Invalid USER_ID format"
// @Failure      400  {object}  models.GetNotificationsResponseSwagger  "Unauthorized access"
// @Failure      500  {object}  map[string]string                 "Internal server error"
// @Security     ApiKeyAuth
// @Router       /v1/notification/list [get]
func (nc *NotificationController) GetNotificationsHandler(c *gin.Context) {
	rawUserID, _ := c.Get("USER_ID")
	userIDStr, ok := rawUserID.(string)
	if !ok {
		c.JSON(http.StatusOK, &models.GetNotificationsResponse{
			Status:            -1,
			StatusDescription: "unable to find USER_ID",
		})
		return
	}

	// Parse to UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusOK, &models.GetNotificationsResponse{
			Status:            -1,
			StatusDescription: "unable to find USER_ID",
		})
		return
	}

	notifications, err := nc.notificationService.GetNotificationsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.GetNotificationsResponse{
			Status:            -1,
			StatusDescription: "Error while fetching Notifications",
		})
		return
	}

	c.JSON(http.StatusOK, models.GetNotificationsResponse{
		Status:            1,
		StatusDescription: "Notifications fetch successful",
		Notifications:     &notifications,
	})
}

// ReadNotificationsHandler godoc
//
// @Summary      Mark notifications as read
// @Description  Marks a list of user notifications as read based on their IDs.
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Param        request  body      models.ReadNotificationsRequest   true  "List of notification IDs to mark as read"
// @Success      200      {object}  models.ReadNotificationsResponse  "Notifications marked as read successfully"
// @Failure      400      {object}  models.ResetPasswordInitResponse  "Invalid request payload"
// @Failure      401      {object}  models.ResetPasswordInitResponse  "Error marking notifications as read"
// @Security     ApiKeyAuth
// @Router       /v1/notification/read/ [post]
func (nc *NotificationController) ReadNotificationsHandler(c *gin.Context) {
	var readNotificationsRequest models.ReadNotificationsRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&readNotificationsRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.ResetPasswordInitResponse{
			Status:            0,
			StatusDescription: "Invalid request payload",
		})
		return
	}

	err := nc.notificationService.MarkNotificationsAsRead(readNotificationsRequest.IDs)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ResetPasswordInitResponse{
			Status:            0,
			StatusDescription: "Error while marking notifications as read",
		})
		return
	}

	c.JSON(http.StatusOK, &models.ReadNotificationsResponse{
		Status:            1,
		StatusDescription: "Notifications marked as read",
	})
}
