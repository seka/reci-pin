package notification

//go:generate mockgen -source=$GOFILE -destination=mock/email_client_mock.go -package=mock
type EmailClient interface {
	SendPasswordChangeNotification(to string) error
	SendPasswordReset(to string, token string) error
}
