package delivery

import (
	"fmt"
	"net/smtp"

	"notification-service/internal/model"
)

type EmailSender struct {
	host     string
	port     string
	username string
	password string
}

func NewEmailSender(host, port, user, pass string) *EmailSender {
	return &EmailSender{
		host:     host,
		port:     port,
		username: user,
		password: pass,
	}
}

func (s *EmailSender) Send(n *model.Notification) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := []byte(fmt.Sprintf(
		"Subject: %s\r\n\r\n%s",
		n.Title,
		n.Body,
	))

	return smtp.SendMail(
		s.host+":"+s.port,
		auth,
		s.username,
		[]string{"user@example.com"}, // todo -> will fetch actual user email
		msg,
	)
}
