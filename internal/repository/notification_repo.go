package repository

import "notification-service/internal/model"

type NotificationRepository interface {
	Create(notification *model.Notification) error
	GetByID(id string) (*model.Notification, error)
	UpdateStatus(id string, status model.NotificationStatus, err *string) error
	IncrementRetry(id string) error
}

// why interface? --> i) easy testing ii) easy to switch db later iii) clean architecture signal
