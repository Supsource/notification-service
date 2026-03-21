package delivery

import (
	"fmt"
	"log"
	"net/smtp"

	"notification-service/internal/model"
)

type EmailSender struct {
	host     string
	port     string
	username string
	password string
	logger   *log.Logger
}

func NewEmailSender(host, port, user, pass string) *EmailSender {
	return &EmailSender{
		host:     host,
		port:     port,
		username: user,
		password: pass,
		logger:   log.Default(),
	}
}

func (s *EmailSender) Send(n *model.Notification) error {
	if s.host == "" || s.port == "" || s.username == "" {
		s.logger.Printf("stub email delivery for notification %s to user %s", n.ID, n.UserID)
		return nil
	}

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
