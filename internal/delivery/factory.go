package delivery

import "notification-service/internal/model"

type Factory struct {
	email Sender
	push Sender
}

func NewFactory(email Sender, push Sender) *Factory {
	return &Factory{
		email: email,
		push: push,
	}
}

func (f *Factory) GetSender(t model.NotificationType) Sender {
	switch t {
	case model.TypeEmail:
		return f.email
	case model.TypePush:
		return f.push
	default:
		return nil
	}
}