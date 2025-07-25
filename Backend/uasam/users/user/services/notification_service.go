package services

import (
	"fmt"
	"uasam/logger"
	"uasam/users/user/models"
	"uasam/users/user/repositories"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NotificationService struct {
	notificationRepository *repositories.NotificationRepository
	logger                 *logger.Logger
}

func NewNotificationService(notificationRepository *repositories.NotificationRepository, logger *logger.Logger) *NotificationService {

	return &NotificationService{
		notificationRepository: notificationRepository,
		logger:                 logger,
	}
}

func (ns *NotificationService) CreateSignUpNotification(user *models.User) error {
	welcomeMessage := fmt.Sprintf("Hi %s, explore zhuhana by creating your first algorithm!", user.FirstName)
	metadata := map[string]interface{}{}
	link := "/code"
	_, err := ns.notificationRepository.CreateNotification(user.ID, "sign_up", "Welcome to Zhuhana", welcomeMessage, &link, metadata)
	if err != nil {
		go ns.logger.Warning("error while creating notification", zap.String("execution level", "CreateSignUpNotification"), zap.String("Error", err.Error()))
	}
	return err
}

func (ns *NotificationService) GetNotificationsByUserID(userID uuid.UUID) ([]models.Notification, error) {
	notifications, err := ns.notificationRepository.GetNotificationsByUserID(userID)
	if err != nil {
		go ns.logger.Warning("error while getting notifications", zap.String("execution level", "GetNotificationsByUserID"), zap.String("Error", err.Error()))
		return nil, err
	}
	return notifications, nil
}

func (ns *NotificationService) MarkNotificationsAsRead(ids []uuid.UUID) error {
	err := ns.notificationRepository.MarkNotificationsAsRead(ids)
	if err != nil {
		go ns.logger.Warning("error while marking notifications as read", zap.String("execution level", "MarkNotificationsAsRead"), zap.String("Error", err.Error()))
	}
	return err
}
