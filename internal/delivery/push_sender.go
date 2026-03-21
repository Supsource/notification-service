package delivery

import (
	"log"

	"notification-service/internal/model"
)

type PushSender struct {
	logger *log.Logger
}

func NewPushSender() *PushSender {
	return &PushSender{logger: log.Default()}
}

func (s *PushSender) Send(n *model.Notification) error {
	s.logger.Printf("stub push delivery for notification %s to user %s", n.ID, n.UserID)
	return nil
}
