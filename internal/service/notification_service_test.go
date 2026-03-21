package service

import (
	"context"
	"testing"
	"time"

	"notification-service/internal/model"
)

type stubNotificationRepo struct{}

func (s *stubNotificationRepo) Create(notification *model.Notification) error  { return nil }
func (s *stubNotificationRepo) GetByID(id string) (*model.Notification, error) { return nil, nil }
func (s *stubNotificationRepo) UpdateStatus(id string, status model.NotificationStatus, err *string) error {
	return nil
}
func (s *stubNotificationRepo) IncrementRetry(id string) error { return nil }

type stubOutboxRepo struct{}

func (s *stubOutboxRepo) Enqueue(ctx context.Context, n *model.OutboxNotification) error {
	return nil
}
func (s *stubOutboxRepo) ClaimPending(ctx context.Context, batchSize int, processingTimeout time.Duration) ([]model.OutboxNotification, error) {
	return nil, nil
}
func (s *stubOutboxRepo) MarkSent(ctx context.Context, id string) error { return nil }
func (s *stubOutboxRepo) ScheduleRetry(ctx context.Context, id string, retryCount int, nextRetryAt time.Time) error {
	return nil
}
func (s *stubOutboxRepo) MarkFailed(ctx context.Context, id string) error { return nil }
func (s *stubOutboxRepo) ListFailed(ctx context.Context, limit, offset int) ([]model.OutboxNotification, error) {
	return nil, nil
}
func (s *stubOutboxRepo) RetryFailed(ctx context.Context, ids []string) (int64, error) { return 0, nil }

func TestCreateNotificationRejectsUnsupportedType(t *testing.T) {
	svc := NewNotificationService(&stubNotificationRepo{}, &stubOutboxRepo{})

	err := svc.CreateNotification("123e4567-e89b-12d3-a456-426614174000", model.NotificationType("SMS"), "hello", "world")
	if err != ErrUnsupportedNotificationType {
		t.Fatalf("expected ErrUnsupportedNotificationType, got %v", err)
	}
}
