package service

import (
	"notification-service/internal/model"
	"notification-service/internal/repository"

	"github.com/google/uuid"
)

type NotificationService struct {
	repo repository.NotificationRepository
}

func (s *NotificationService) CreateNotification(
	userID string,
	nType model.NotificationType,
	title string,
	body string,
) error {
	notification := &model.Notification{
		ID:         uuid.New().String(),
		UserID:     userID,
		Type:       nType,
		Title:      title,
		Body:       body,
		Status:     model.StatusPending,
		RetryCount: 0,
	}

	return s.repo.Create(notification)
}

// service generates UUID, sets status and repo just saves
