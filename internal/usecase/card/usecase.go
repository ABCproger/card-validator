package card

import (
	"context"
	"time"

	"github.com/ABCproger/card-validator/internal/entity"
)

type UseCase struct {
	now func() time.Time
}

func New() *UseCase {
	return &UseCase{now: time.Now}
}

func NewWithClock(now func() time.Time) *UseCase {
	return &UseCase{now: now}
}

func (uc *UseCase) Validate(_ context.Context, card entity.Card) *entity.ValidationError {
	return card.Validate(uc.now())
}
