package model

import "time"

type NotificationStatus string
type NotificationType string

const (
	StatusPending    NotificationStatus = "PENDING"
	StatusProcessing NotificationStatus = "PROCESSING"
	StatusSent       NotificationStatus = "SENT"
	StatusFailed     NotificationStatus = "FAILED"

	TypeInApp NotificationType = "IN_APP"
	TypeEmail NotificationType = "EMAIL"
	TypePush  NotificationType = "PUSH"
)

type Notification struct {
	ID          string
	UserID      string
	Type        NotificationType
	Title       string
	Body        string
	Status      NotificationStatus
	RetryCount  int
	ErrorReason *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
