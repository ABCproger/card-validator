package card

import (
	"context"
	"testing"
	"time"

	"github.com/ABCproger/card-validator/internal/entity"
	"github.com/stretchr/testify/assert"
)

func fixedNow() time.Time {
	return time.Date(2026, 6, 16, 12, 0, 0, 0, time.UTC)
}

func TestUseCase_Validate(t *testing.T) {
	uc := NewWithClock(fixedNow)
	ctx := context.Background()

	tests := []struct {
		name      string
		card      entity.Card
		wantCode  string
		wantValid bool
	}{
		{name: "valid visa", card: entity.Card{Number: "4111111111111111", ExpirationMonth: 12, ExpirationYear: 2028}, wantValid: true},
		{name: "invalid number", card: entity.Card{Number: "1111111111111", ExpirationMonth: 12, ExpirationYear: 2028}, wantCode: entity.ErrCodeInvalidNumber},
		{name: "invalid month", card: entity.Card{Number: "4111111111111111", ExpirationMonth: 13, ExpirationYear: 2028}, wantCode: entity.ErrCodeInvalidMonth},
		{name: "invalid year", card: entity.Card{Number: "4111111111111111", ExpirationMonth: 12, ExpirationYear: 1900}, wantCode: entity.ErrCodeExpired},
		{name: "expired", card: entity.Card{Number: "4111111111111111", ExpirationMonth: 1, ExpirationYear: 2021}, wantCode: entity.ErrCodeExpired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verr := uc.Validate(ctx, tt.card)
			if tt.wantValid {
				assert.Nil(t, verr)
			} else {
				assert.NotNil(t, verr)
				assert.Equal(t, tt.wantCode, verr.Code)
			}
		})
	}
}
