package delivery

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"

	"notification-service/internal/model"
)

type PushSender struct {
	client *messaging.Client
}

func NewPushSender(app *firebase.App) (*PushSender, error) {
	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}
	return &PushSender{client: client}, nil
}

func (s *PushSender) Send(n *model.Notification) error {
	msg := &messaging.Message{
		Notification: &messaging.Notification{
			Title: n.Title,
			Body:  n.Body,
		},
		Token: "DEVICE_FCM_TOCKEN", // todo -> later from db
	}

	_, err := s.client.Send(context.Background(), msg)
	return err
}
