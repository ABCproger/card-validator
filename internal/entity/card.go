package entity

import (
	"strings"
	"time"
)

const (
	ErrCodeInvalidNumber = "001"
	ErrCodeInvalidMonth  = "002"
	ErrCodeExpired       = "003"
)

type ValidationError struct {
	Code    string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type Card struct {
	Number          string
	ExpirationMonth int
	ExpirationYear  int
}

func (c Card) Validate(now time.Time) *ValidationError {
	if err := c.validateNumber(); err != nil {
		return err
	}
	if err := c.validateMonth(); err != nil {
		return err
	}
	if err := c.validateExpiry(now); err != nil {
		return err
	}
	return nil
}

func (c Card) validateNumber() *ValidationError {
	normalized := normalize(c.Number)
	if !allASCIIDigits(normalized) {
		return &ValidationError{Code: ErrCodeInvalidNumber, Message: "card number must contain digits only"}
	}
	if len(normalized) < 8 || len(normalized) > 19 {
		return &ValidationError{Code: ErrCodeInvalidNumber, Message: "card number length must be between 8 and 19 digits"}
	}
	if !luhnValid(normalized) {
		return &ValidationError{Code: ErrCodeInvalidNumber, Message: "card number is invalid (failed Luhn check)"}
	}
	return nil
}

func (c Card) validateMonth() *ValidationError {
	if c.ExpirationMonth < 1 || c.ExpirationMonth > 12 {
		return &ValidationError{Code: ErrCodeInvalidMonth, Message: "expiration month must be between 1 and 12"}
	}
	return nil
}

func (c Card) validateExpiry(now time.Time) *ValidationError {
	currentYear, currentMonth, _ := now.Date()

	expired := c.ExpirationYear < currentYear ||
		(c.ExpirationYear == currentYear && c.ExpirationMonth < int(currentMonth))

	if expired {
		return &ValidationError{Code: ErrCodeExpired, Message: "card has expired"}
	}
	return nil
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	return s
}

func allASCIIDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func luhnValid(number string) bool {
	sum := 0
	double := false
	for i := len(number) - 1; i >= 0; i-- {
		d := int(number[i] - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return sum%10 == 0
}
