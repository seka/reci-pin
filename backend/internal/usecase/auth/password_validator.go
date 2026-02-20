package auth

import (
	"regexp"

	"github.com/seka/reci-pin/backend/internal/domain/validation"
)

func ValidatePassword(password string) error {
	var errs validation.ValidationErrors

	if len(password) < 8 {
		errs = append(errs, validation.ValidationError{
			Field:  "password",
			Code:   validation.ErrCodePasswordTooShort,
			Params: map[string]any{"min": 8},
		})
	}

	hasAlpha := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	if !hasAlpha {
		errs = append(errs, validation.ValidationError{
			Field: "password",
			Code:  validation.ErrCodePasswordNoAlpha,
		})
	}

	hasNumeric := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumeric {
		errs = append(errs, validation.ValidationError{
			Field: "password",
			Code:  validation.ErrCodePasswordNoNumeric,
		})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
