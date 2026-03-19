package model

import (
	"encoding/json"
	"time"
)

type OutboxStatus string

const (
	OutboxStatusPending    OutboxStatus = "pending"
	OutboxStatusProcessing OutboxStatus = "processing"
	OutboxStatusSent       OutboxStatus = "sent"
	OutboxStatusFailed     OutboxStatus = "failed"
)

type OutboxNotification struct {
	ID          string
	UserID      string
	Payload     json.RawMessage
	Status      OutboxStatus
	RetryCount  int
	NextRetryAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OutboxPayload struct {
	NotificationID string           `json:"notification_id"`
	UserID         string           `json:"user_id"`
	Type           NotificationType `json:"type"`
	Title          string           `json:"title"`
	Body           string           `json:"body"`
}
