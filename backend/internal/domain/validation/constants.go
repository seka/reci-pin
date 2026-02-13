package validation

// Validation Constants
const (
	EmailMaxLength      = 200
	PasswordMinLength   = 8
	PasswordMaxLength   = 200
	RecipeNameMaxLength = 100
	TagNameMaxLength    = 30
)

// Common Error Codes
const (
	ErrCodeEmailInvalidFormat = "EMAIL_INVALID_FORMAT"
	ErrCodeEmailTooLong       = "EMAIL_TOO_LONG"
	ErrCodeRequired           = "REQUIRED"
	ErrCodeTextTooLong        = "TEXT_TOO_LONG"
	ErrCodeURLInvalidFormat   = "URL_INVALID_FORMAT"
	ErrCodePasswordTooShort   = "PASSWORD_TOO_SHORT"
	ErrCodePasswordTooLong    = "PASSWORD_TOO_LONG"
	ErrCodePasswordNoAlpha    = "PASSWORD_NO_ALPHA"
	ErrCodePasswordNoNumeric  = "PASSWORD_NO_NUMERIC"
)
