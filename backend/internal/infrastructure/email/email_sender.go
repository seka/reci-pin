package email

import (
	"fmt"
	"net/smtp"
)

type EmailSender interface {
	SendPasswordChangeNotification(to string) error
}

type MailHogSender struct {
	host string
	port int
	from string
}

func NewMailHogSender(host string, port int, from string) *MailHogSender {
	return &MailHogSender{
		host: host,
		port: port,
		from: from,
	}
}

func (s *MailHogSender) SendPasswordChangeNotification(to string) error {
	subject := "Password Changed Successfully"
	body := "Your password has been changed successfully.\nIf you did not make this change, please contact support immediately."

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.from, to, subject, body)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	// MailHog uses no authentication by default
	if err := smtp.SendMail(addr, nil, s.from, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
