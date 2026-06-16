package v1

import (
	"google.golang.org/protobuf/proto"

	cardv1 "github.com/ABCproger/card-validator/gen/card/v1"
	"github.com/ABCproger/card-validator/internal/entity"
)

func toProtoResponse(verr *entity.ValidationError) *cardv1.ValidateResponse {
	resp := &cardv1.ValidateResponse{Valid: proto.Bool(verr == nil)}
	if verr != nil {
		resp.Error = &cardv1.ValidationError{
			Code:    verr.Code,
			Message: verr.Message,
		}
	}
	return resp
}
