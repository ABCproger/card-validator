package card

import (
	"context"

	"github.com/ABCproger/card-validator/internal/entity"
)

type Validator interface {
	Validate(ctx context.Context, card entity.Card) *entity.ValidationError
}
