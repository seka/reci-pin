package validation

// Validation Constants
const (
	EmailMaxLength      = 200
	PasswordMinLength   = 8
	PasswordMaxLength   = 200
	RecipeNameMaxLength = 100
	TagNameMaxLength    = 30

	// 最近のスマホ(ProRAW/高画素モード等)では10MBを超える画像が生成されるため、余裕を持たせて50MBに設定
	ImageMaxFileSize = 50 * 1024 * 1024 // 50MB
)

// Allowed image content types for upload
var ImageAllowedContentTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

// Allowed image file extensions for upload
var ImageAllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

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
