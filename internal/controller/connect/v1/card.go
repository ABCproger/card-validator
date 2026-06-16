package v1

import (
	"context"

	"connectrpc.com/connect"

	cardv1 "github.com/ABCproger/card-validator/gen/card/v1"
	"github.com/ABCproger/card-validator/gen/card/v1/cardv1connect"
	"github.com/ABCproger/card-validator/internal/entity"
	"github.com/ABCproger/card-validator/internal/usecase/card"
)

type CardHandler struct {
	cardv1connect.UnimplementedCardServiceHandler
	validator card.Validator
}

func NewCardHandler(v card.Validator) *CardHandler {
	return &CardHandler{validator: v}
}

func (h *CardHandler) Validate(
	ctx context.Context,
	req *connect.Request[cardv1.ValidateRequest],
) (*connect.Response[cardv1.ValidateResponse], error) {
	c := entity.Card{
		Number:          req.Msg.CardNumber,
		ExpirationMonth: int(req.Msg.ExpirationMonth),
		ExpirationYear:  int(req.Msg.ExpirationYear),
	}

	verr := h.validator.Validate(ctx, c)

	return connect.NewResponse(toProtoResponse(verr)), nil
}
