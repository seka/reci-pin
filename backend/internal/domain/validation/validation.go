package validation

import (
	"fmt"
	"net/url"
	"regexp"
	"unicode/utf8"
)

type ValidationError struct {
	Field  string
	Code   string
	Params map[string]interface{}
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Code)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	return "validation failed"
}

// ValidateEmail checks if the email format is valid and length is within limits.
func ValidateEmail(email string) error {
	var errs ValidationErrors

	if email == "" {
		errs = append(errs, ValidationError{Field: "email", Code: ErrCodeRequired})
		return errs
	}

	if len(email) > EmailMaxLength {
		errs = append(errs, ValidationError{
			Field:  "email",
			Code:   ErrCodeEmailTooLong,
			Params: map[string]interface{}{"max": EmailMaxLength},
		})
	}

	// Simple regex for email validation
	// Matches: something@something.something
	emailRegex := regexp.MustCompile(`^[\w+\-.]+@[a-z\d\-]+(\.[a-z\d\-]+)*\.[a-z]+$`)
	if !emailRegex.MatchString(email) {
		errs = append(errs, ValidationError{
			Field: "email",
			Code:  ErrCodeEmailInvalidFormat,
		})
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// ValidatePassword checks password complexity.
func ValidatePassword(password string) error {
	var errs ValidationErrors

	if len(password) < PasswordMinLength {
		errs = append(errs, ValidationError{
			Field:  "password",
			Code:   ErrCodePasswordTooShort,
			Params: map[string]interface{}{"min": PasswordMinLength},
		})
	}
	if len(password) > PasswordMaxLength {
		errs = append(errs, ValidationError{
			Field:  "password",
			Code:   ErrCodePasswordTooLong,
			Params: map[string]interface{}{"max": PasswordMaxLength},
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

// ValidateRecipe checks recipe name length and URL format.
func ValidateRecipe(name, link string) error {
	var errs ValidationErrors

	if name == "" {
		errs = append(errs, ValidationError{Field: "name", Code: ErrCodeRequired})
	} else if utf8.RuneCountInString(name) > RecipeNameMaxLength {
		errs = append(errs, ValidationError{
			Field:  "name",
			Code:   ErrCodeTextTooLong,
			Params: map[string]interface{}{"max": RecipeNameMaxLength},
		})
	}

	if link != "" {
		u, err := url.ParseRequestURI(link)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			errs = append(errs, ValidationError{
				Field: "url",
				Code:  ErrCodeURLInvalidFormat,
			})
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// ValidateTag checks tag name length.
func ValidateTag(name string) error {
	var errs ValidationErrors

	if name == "" {
		errs = append(errs, ValidationError{Field: "name", Code: ErrCodeRequired})
	} else if utf8.RuneCountInString(name) > TagNameMaxLength {
		errs = append(errs, ValidationError{
			Field:  "name",
			Code:   ErrCodeTextTooLong,
			Params: map[string]interface{}{"max": TagNameMaxLength},
		})
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
