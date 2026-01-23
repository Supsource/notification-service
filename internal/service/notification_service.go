package service

import (
	"notification-service/internal/model"
	"notification-service/internal/queue"
	"notification-service/internal/repository"

	"github.com/google/uuid"
)

type NotificationService struct {
	repo     repository.NotificationRepository
	producer *queue.Producer
}

func NewNotificationService(repo repository.NotificationRepository, producer *queue.Producer) *NotificationService {
	return &NotificationService{
		repo:     repo,
		producer: producer,
	}
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

	if err := s.repo.Create(notification); err != nil {
		return err
	}

	return s.producer.Enqueue(notification.ID)
}
