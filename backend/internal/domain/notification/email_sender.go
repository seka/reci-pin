package notification

//go:generate mockgen -source=$GOFILE -destination=mock/email_sender_mock.go -package=mock
type EmailSender interface {
	SendPasswordChangeNotification(to string) error
	SendPasswordReset(to string, token string) error
}
