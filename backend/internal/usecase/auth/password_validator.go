package auth

import (
	"regexp"
)

const (
	ErrCodePasswordTooShort  = "PASSWORD_TOO_SHORT"
	ErrCodePasswordNoAlpha   = "PASSWORD_NO_ALPHA"
	ErrCodePasswordNoNumeric = "PASSWORD_NO_NUMERIC"
)

type ValidationError struct {
	Field  string
	Code   string
	Params map[string]interface{}
}

func (e ValidationError) Error() string {
	return e.Code
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	return "validation failed"
}

func ValidatePassword(password string) error {
	var errs ValidationErrors

	if len(password) < 8 {
		errs = append(errs, ValidationError{
			Field:  "password",
			Code:   ErrCodePasswordTooShort,
			Params: map[string]interface{}{"min": 8},
		})
	}

	hasAlpha := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	if !hasAlpha {
		errs = append(errs, ValidationError{
			Field: "password",
			Code:  ErrCodePasswordNoAlpha,
		})
	}

	hasNumeric := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumeric {
		errs = append(errs, ValidationError{
			Field: "password",
			Code:  ErrCodePasswordNoNumeric,
		})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
