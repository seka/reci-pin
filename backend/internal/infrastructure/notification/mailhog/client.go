package mailhog

import (
	"fmt"
	"net/smtp"

	"github.com/seka/reci-pin/backend/internal/domain/notification"
)

type Client struct {
	host string
	port int
	from string
}

func New(host string, port int, from string) notification.EmailSender {
	return &Client{
		host: host,
		port: port,
		from: from,
	}
}

func (s *Client) SendPasswordChangeNotification(to string) error {
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

func (s *Client) SendPasswordReset(to string, token string) error {
	resetURL := fmt.Sprintf("http://localhost:4200/password-reset?token=%s", token)
	subject := "パスワード再設定のご案内"
	body := fmt.Sprintf("いつも Reci-pin をご利用いただきありがとうございます。\n"+
		"パスワード再設定のリクエストを受け付けました。\n\n"+
		"以下のリンクをクリックして、新しいパスワードを設定してください。\n\n"+
		"%s\n\n"+
		"※このリンクは30分間有効です。\n"+
		"※本メールに心当たりがない場合は、破棄してください。\n\n"+
		"--------------------------------------------------\n"+
		"Reci-pin 運営事務局\n"+
		"お問い合わせ: support@reci-pin.com\n"+
		"プライバシーポリシー: https://reci-pin.com/privacy\n"+
		"--------------------------------------------------", resetURL)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.from, to, subject, body)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	if err := smtp.SendMail(addr, nil, s.from, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
