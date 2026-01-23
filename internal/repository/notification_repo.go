package repository

import "notification-service/internal/model"

type NotificationRepository interface {
	Create(notification *model.Notification) error
}

// why interface? --> i) easy testing ii) easy to switch db later iii) clean architecture signal
