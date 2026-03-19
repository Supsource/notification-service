package workers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"notification-service/internal/delivery"
	"notification-service/internal/model"
	"notification-service/internal/repository"
)

const (
	defaultBatchSize         = 50
	defaultPollInterval      = 2 * time.Second
	defaultProcessingTimeout = 1 * time.Minute
	maxRetries               = 5
)

type OutboxWorker struct {
	outboxRepo        repository.OutboxRepository
	factory           *delivery.Factory
	batchSize         int
	pollInterval      time.Duration
	processingTimeout time.Duration
	logger            *log.Logger
}

func NewOutboxWorker(outboxRepo repository.OutboxRepository, factory *delivery.Factory, logger *log.Logger) *OutboxWorker {
	if logger == nil {
		logger = log.Default()
	}
	return &OutboxWorker{
		outboxRepo:        outboxRepo,
		factory:           factory,
		batchSize:         defaultBatchSize,
		pollInterval:      defaultPollInterval,
		processingTimeout: defaultProcessingTimeout,
		logger:            logger,
	}
}

func (w *OutboxWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		processed, err := w.processBatch(ctx)
		if err != nil {
			w.logger.Println("outbox worker error:", err)
		}
		if processed == 0 {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}
}

func (w *OutboxWorker) processBatch(ctx context.Context) (int, error) {
	items, err := w.outboxRepo.ClaimPending(ctx, w.batchSize, w.processingTimeout)
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		w.processItem(ctx, item)
	}
	return len(items), nil
}

func (w *OutboxWorker) processItem(ctx context.Context, item model.OutboxNotification) {
	var payload model.OutboxPayload
	if err := json.Unmarshal(item.Payload, &payload); err != nil {
		w.logger.Printf("outbox %s invalid payload: %v", item.ID, err)
		_ = w.outboxRepo.MarkFailed(ctx, item.ID)
		return
	}

	notification := &model.Notification{
		ID:     payload.NotificationID,
		UserID: payload.UserID,
		Type:   payload.Type,
		Title:  payload.Title,
		Body:   payload.Body,
	}

	sender := w.factory.GetSender(notification.Type)
	if sender == nil {
		w.logger.Printf("outbox %s unsupported type: %s", item.ID, notification.Type)
		_ = w.outboxRepo.MarkFailed(ctx, item.ID)
		return
	}

	if err := sender.Send(notification); err != nil {
		retryCount := item.RetryCount + 1
		if retryCount > maxRetries {
			w.logger.Printf("outbox %s failed permanently after %d retries: %v", item.ID, item.RetryCount, err)
			_ = w.outboxRepo.MarkFailed(ctx, item.ID)
			return
		}

		nextRetry := time.Now().Add(backoffDuration(retryCount))
		w.logger.Printf("outbox %s retry %d scheduled at %s: %v", item.ID, retryCount, nextRetry.Format(time.RFC3339), err)
		_ = w.outboxRepo.ScheduleRetry(ctx, item.ID, retryCount, nextRetry)
		return
	}

	w.logger.Printf("outbox %s sent successfully", item.ID)
	_ = w.outboxRepo.MarkSent(ctx, item.ID)
}

func backoffDuration(retryCount int) time.Duration {
	switch retryCount {
	case 1:
		return 10 * time.Second
	case 2:
		return 30 * time.Second
	case 3:
		return 2 * time.Minute
	case 4:
		return 10 * time.Minute
	default:
		return 10 * time.Minute
	}
}
