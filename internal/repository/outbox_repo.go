package repository

import (
	"context"
	"time"

	"notification-service/internal/model"
)

type OutboxRepository interface {
	Enqueue(ctx context.Context, n *model.OutboxNotification) error
	ClaimPending(ctx context.Context, batchSize int, processingTimeout time.Duration) ([]model.OutboxNotification, error)
	MarkSent(ctx context.Context, id string) error
	ScheduleRetry(ctx context.Context, id string, retryCount int, nextRetryAt time.Time) error
	MarkFailed(ctx context.Context, id string) error
	ListFailed(ctx context.Context, limit, offset int) ([]model.OutboxNotification, error)
	RetryFailed(ctx context.Context, ids []string) (int64, error)
}
