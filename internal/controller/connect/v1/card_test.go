package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cardv1 "github.com/ABCproger/card-validator/gen/card/v1"
	"github.com/ABCproger/card-validator/gen/card/v1/cardv1connect"
	"github.com/ABCproger/card-validator/internal/entity"
	uccard "github.com/ABCproger/card-validator/internal/usecase/card"
)

func newTestServer() *httptest.Server {
	uc := uccard.NewWithClock(func() time.Time {
		return time.Date(2026, 6, 16, 12, 0, 0, 0, time.UTC)
	})
	handler := NewCardHandler(uc)
	mux := http.NewServeMux()
	mux.Handle(cardv1connect.NewCardServiceHandler(handler))
	return httptest.NewUnstartedServer(mux)
}

func TestCardHandler_Validate(t *testing.T) {
	srv := newTestServer()
	srv.EnableHTTP2 = false
	srv.Start()
	defer srv.Close()

	client := cardv1connect.NewCardServiceClient(srv.Client(), srv.URL)

	tests := []struct {
		name      string
		req       *cardv1.ValidateRequest
		wantValid bool
		wantCode  string
	}{
		{
			name:      "valid visa",
			req:       &cardv1.ValidateRequest{CardNumber: "4111111111111111", ExpirationMonth: 12, ExpirationYear: 2028},
			wantValid: true,
		},
		{
			name:     "invalid number (Luhn fail)",
			req:      &cardv1.ValidateRequest{CardNumber: "1111111111111", ExpirationMonth: 10, ExpirationYear: 2028},
			wantCode: entity.ErrCodeInvalidNumber,
		},
		{
			name:     "expired card",
			req:      &cardv1.ValidateRequest{CardNumber: "4111111111111111", ExpirationMonth: 1, ExpirationYear: 2021},
			wantCode: entity.ErrCodeExpired,
		},
		{
			name:     "invalid month",
			req:      &cardv1.ValidateRequest{CardNumber: "4111111111111111", ExpirationMonth: 13, ExpirationYear: 2028},
			wantCode: entity.ErrCodeInvalidMonth,
		},
		{
			name:      "card valid through end of current month",
			req:       &cardv1.ValidateRequest{CardNumber: "4111111111111111", ExpirationMonth: 6, ExpirationYear: 2026},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Validate(context.Background(), connect.NewRequest(tt.req))
			require.NoError(t, err)

			if tt.wantValid {
				assert.True(t, resp.Msg.GetValid())
				assert.Nil(t, resp.Msg.Error)
			} else {
				assert.False(t, resp.Msg.GetValid())
				require.NotNil(t, resp.Msg.Error)
				assert.Equal(t, tt.wantCode, resp.Msg.Error.Code)
			}
		})
	}
}
