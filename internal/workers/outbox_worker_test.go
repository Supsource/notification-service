package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"testing"
	"time"

	"notification-service/internal/delivery"
	"notification-service/internal/model"
)

type stubNotificationRepo struct {
	statusUpdates []statusUpdate
	retryIDs      []string
}

type statusUpdate struct {
	id     string
	status model.NotificationStatus
	err    *string
}

func (s *stubNotificationRepo) Create(notification *model.Notification) error  { return nil }
func (s *stubNotificationRepo) GetByID(id string) (*model.Notification, error) { return nil, nil }
func (s *stubNotificationRepo) UpdateStatus(id string, status model.NotificationStatus, err *string) error {
	s.statusUpdates = append(s.statusUpdates, statusUpdate{id: id, status: status, err: err})
	return nil
}
func (s *stubNotificationRepo) IncrementRetry(id string) error {
	s.retryIDs = append(s.retryIDs, id)
	return nil
}

type stubOutboxRepo struct {
	markedSent   []string
	markedFailed []string
	scheduled    []scheduledRetry
}

type scheduledRetry struct {
	id         string
	retryCount int
	nextRetry  time.Time
}

func (s *stubOutboxRepo) Enqueue(ctx context.Context, n *model.OutboxNotification) error { return nil }
func (s *stubOutboxRepo) ClaimPending(ctx context.Context, batchSize int, processingTimeout time.Duration) ([]model.OutboxNotification, error) {
	return nil, nil
}
func (s *stubOutboxRepo) MarkSent(ctx context.Context, id string) error {
	s.markedSent = append(s.markedSent, id)
	return nil
}
func (s *stubOutboxRepo) ScheduleRetry(ctx context.Context, id string, retryCount int, nextRetryAt time.Time) error {
	s.scheduled = append(s.scheduled, scheduledRetry{id: id, retryCount: retryCount, nextRetry: nextRetryAt})
	return nil
}
func (s *stubOutboxRepo) MarkFailed(ctx context.Context, id string) error {
	s.markedFailed = append(s.markedFailed, id)
	return nil
}
func (s *stubOutboxRepo) ListFailed(ctx context.Context, limit, offset int) ([]model.OutboxNotification, error) {
	return nil, nil
}
func (s *stubOutboxRepo) RetryFailed(ctx context.Context, ids []string) (int64, error) { return 0, nil }

type stubSender struct {
	err error
}

func (s *stubSender) Send(n *model.Notification) error {
	return s.err
}

func TestProcessItemMarksNotificationSentOnSuccess(t *testing.T) {
	notificationRepo := &stubNotificationRepo{}
	outboxRepo := &stubOutboxRepo{}
	worker := NewOutboxWorker(
		notificationRepo,
		outboxRepo,
		delivery.NewFactory(&stubSender{}, &stubSender{}),
		log.New(io.Discard, "", 0),
	)

	payload, _ := json.Marshal(model.OutboxPayload{
		NotificationID: "notif-1",
		UserID:         "user-1",
		Type:           model.TypeEmail,
		Title:          "hello",
		Body:           "world",
	})

	worker.processItem(context.Background(), model.OutboxNotification{
		ID:      "outbox-1",
		Payload: payload,
	})

	if len(notificationRepo.statusUpdates) != 2 {
		t.Fatalf("expected 2 status updates, got %d", len(notificationRepo.statusUpdates))
	}
	if notificationRepo.statusUpdates[0].status != model.StatusProcessing {
		t.Fatalf("expected first status to be processing, got %s", notificationRepo.statusUpdates[0].status)
	}
	if notificationRepo.statusUpdates[1].status != model.StatusSent {
		t.Fatalf("expected final status to be sent, got %s", notificationRepo.statusUpdates[1].status)
	}
	if len(outboxRepo.markedSent) != 1 || outboxRepo.markedSent[0] != "outbox-1" {
		t.Fatalf("expected outbox item to be marked sent, got %#v", outboxRepo.markedSent)
	}
}

func TestProcessItemSchedulesRetryAndResetsNotificationStatus(t *testing.T) {
	notificationRepo := &stubNotificationRepo{}
	outboxRepo := &stubOutboxRepo{}
	worker := NewOutboxWorker(
		notificationRepo,
		outboxRepo,
		delivery.NewFactory(&stubSender{err: errors.New("temporary failure")}, &stubSender{}),
		log.New(io.Discard, "", 0),
	)

	payload, _ := json.Marshal(model.OutboxPayload{
		NotificationID: "notif-2",
		UserID:         "user-2",
		Type:           model.TypeEmail,
		Title:          "hello",
		Body:           "world",
	})

	worker.processItem(context.Background(), model.OutboxNotification{
		ID:         "outbox-2",
		Payload:    payload,
		RetryCount: 1,
	})

	if len(notificationRepo.retryIDs) != 1 || notificationRepo.retryIDs[0] != "notif-2" {
		t.Fatalf("expected notification retry increment, got %#v", notificationRepo.retryIDs)
	}
	if len(outboxRepo.scheduled) != 1 {
		t.Fatalf("expected one scheduled retry, got %d", len(outboxRepo.scheduled))
	}
	if outboxRepo.scheduled[0].retryCount != 2 {
		t.Fatalf("expected retry count 2, got %d", outboxRepo.scheduled[0].retryCount)
	}
	if len(notificationRepo.statusUpdates) != 2 || notificationRepo.statusUpdates[1].status != model.StatusPending {
		t.Fatalf("expected notification to return to pending, got %#v", notificationRepo.statusUpdates)
	}
}

func TestProcessItemMarksInvalidPayloadFailed(t *testing.T) {
	notificationRepo := &stubNotificationRepo{}
	outboxRepo := &stubOutboxRepo{}
	worker := NewOutboxWorker(
		notificationRepo,
		outboxRepo,
		delivery.NewFactory(&stubSender{}, &stubSender{}),
		log.New(bytes.NewBuffer(nil), "", 0),
	)

	worker.processItem(context.Background(), model.OutboxNotification{
		ID:      "outbox-3",
		Payload: []byte("{not-json"),
	})

	if len(outboxRepo.markedFailed) != 1 || outboxRepo.markedFailed[0] != "outbox-3" {
		t.Fatalf("expected outbox item to be marked failed, got %#v", outboxRepo.markedFailed)
	}
	if len(notificationRepo.statusUpdates) != 0 {
		t.Fatalf("expected no notification status updates for invalid payload, got %#v", notificationRepo.statusUpdates)
	}
}
