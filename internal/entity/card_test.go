package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var now = time.Date(2026, 6, 16, 12, 0, 0, 0, time.UTC)

func TestCard_Validate(t *testing.T) {
	tests := []struct {
		name      string
		card      Card
		wantCode  string
		wantValid bool
	}{
		// Valid cards
		{name: "visa 16-digit", card: Card{"4111111111111111", 12, 2028}, wantValid: true},
		{name: "visa 16-digit alt", card: Card{"4012888888881881", 12, 2028}, wantValid: true},
		{name: "visa 13-digit", card: Card{"4222222222222", 12, 2028}, wantValid: true},
		{name: "mastercard 5xxx", card: Card{"5555555555554444", 12, 2028}, wantValid: true},
		{name: "mastercard 5xxx alt", card: Card{"5105105105105100", 12, 2028}, wantValid: true},
		{name: "mastercard 2xxx", card: Card{"2221000000000009", 12, 2028}, wantValid: true},
		{name: "amex 15-digit", card: Card{"378282246310005", 12, 2028}, wantValid: true},
		{name: "amex 15-digit alt", card: Card{"371449635398431", 12, 2028}, wantValid: true},
		{name: "discover", card: Card{"6011111111111117", 12, 2028}, wantValid: true},
		{name: "spaces normalized", card: Card{"4111 1111 1111 1111", 12, 2028}, wantValid: true},
		{name: "dashes normalized", card: Card{"4111-1111-1111-1111", 12, 2028}, wantValid: true},
		{name: "expires current month", card: Card{"4111111111111111", 6, 2026}, wantValid: true},

		// Invalid: card number (001)
		{name: "luhn fail", card: Card{"1111111111111", 10, 2028}, wantCode: ErrCodeInvalidNumber},
		{name: "luhn fail 16-digit", card: Card{"4111111111111112", 12, 2028}, wantCode: ErrCodeInvalidNumber},
		{name: "too short", card: Card{"411111111111", 12, 2028}, wantCode: ErrCodeInvalidNumber},
		{name: "too long", card: Card{"41111111111111111111", 12, 2028}, wantCode: ErrCodeInvalidNumber},
		{name: "non-digits", card: Card{"4111111111111abc", 12, 2028}, wantCode: ErrCodeInvalidNumber},
		{name: "empty number", card: Card{"", 12, 2028}, wantCode: ErrCodeInvalidNumber},

		// Invalid: month (002)
		{name: "month 0", card: Card{"4111111111111111", 0, 2028}, wantCode: ErrCodeInvalidMonth},
		{name: "month 13", card: Card{"4111111111111111", 13, 2028}, wantCode: ErrCodeInvalidMonth},

		// Invalid: expired (003)
		{name: "year 1900", card: Card{"4111111111111111", 12, 1900}, wantCode: ErrCodeExpired},
		{name: "expired jan 2021", card: Card{"4111111111111111", 1, 2021}, wantCode: ErrCodeExpired},
		{name: "expired may 2026", card: Card{"4111111111111111", 5, 2026}, wantCode: ErrCodeExpired},

		// First-failure ordering: invalid number before expired
		{name: "bad number + expired, code 001 first", card: Card{"1111111111111", 1, 2021}, wantCode: ErrCodeInvalidNumber},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.card.Validate(now)
			if tt.wantValid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantCode, err.Code)
			}
		})
	}
}

func TestLuhnValid(t *testing.T) {
	valid := []string{
		"4111111111111111",
		"4012888888881881",
		"4222222222222",
		"5555555555554444",
		"378282246310005",
		"6011111111111117",
	}
	for _, n := range valid {
		assert.True(t, luhnValid(n), "expected Luhn valid: %s", n)
	}

	invalid := []string{
		"1111111111111",
		"4111111111111112",
		"1234567890123456",
	}
	for _, n := range invalid {
		assert.False(t, luhnValid(n), "expected Luhn invalid: %s", n)
	}
}
