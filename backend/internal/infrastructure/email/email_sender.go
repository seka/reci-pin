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
	subject := "パスワード変更完了のお知らせ"
	body := "いつも Reci-pin をご利用いただきありがとうございます。\n" +
		"パスワードの変更が完了しました。\n\n" +
		"もし、この変更に心当たりがない場合は、速やかに運営事務局までお問い合わせください。\n\n" +
		"--------------------------------------------------\n" +
		"Reci-pin 運営事務局\n" +
		"お問い合わせ: support@reci-pin.com\n" +
		"プライバシーポリシー: https://reci-pin.com/privacy\n" +
		"--------------------------------------------------"

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
