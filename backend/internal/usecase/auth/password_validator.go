package auth

import (
	"errors"
	"regexp"
)

var (
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters long")
	ErrPasswordNoAlpha    = errors.New("password must contain at least one letter")
	ErrPasswordNoNumeric  = errors.New("password must contain at least one number")
	ErrPasswordComplexity = errors.New("password must contain at least one letter and one number")
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	hasAlpha := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumeric := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasAlpha || !hasNumeric {
		return ErrPasswordComplexity
	}

	return nil
}
