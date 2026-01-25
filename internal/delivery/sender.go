package delivery

import "notification-service/internal/model"

type Sender interface {
	Send(n *model.Notification) error
}