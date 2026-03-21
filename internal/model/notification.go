package model

import (
	"fmt"
	"strings"
	"time"
)

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

func ParseNotificationType(raw string) (NotificationType, error) {
	notificationType := NotificationType(strings.ToUpper(strings.TrimSpace(raw)))
	if !notificationType.IsSupported() {
		return "", fmt.Errorf("unsupported notification type: %s", raw)
	}
	return notificationType, nil
}

func (t NotificationType) IsSupported() bool {
	switch t {
	case TypeEmail, TypePush:
		return true
	default:
		return false
	}
}

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
