package service

import (
	"context"
	"encoding/json"
	"time"

	"notification-service/internal/model"
	"notification-service/internal/repository"

	"github.com/google/uuid"
)

type NotificationService struct {
	repo       repository.NotificationRepository
	outboxRepo repository.OutboxRepository
}

func NewNotificationService(repo repository.NotificationRepository, outboxRepo repository.OutboxRepository) *NotificationService {
	return &NotificationService{
		repo:       repo,
		outboxRepo: outboxRepo,
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

	payload := model.OutboxPayload{
		NotificationID: notification.ID,
		UserID:         userID,
		Type:           nType,
		Title:          title,
		Body:           body,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	outbox := &model.OutboxNotification{
		ID:          uuid.New().String(),
		UserID:      userID,
		Payload:     raw,
		Status:      model.OutboxStatusPending,
		RetryCount:  0,
		NextRetryAt: time.Now(),
	}

	return s.outboxRepo.Enqueue(context.Background(), outbox)
}
