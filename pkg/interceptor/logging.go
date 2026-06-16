package interceptor

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
)

func NewLoggingInterceptor(log *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			requestID := uuid.New().String()
			start := time.Now()

			resp, err := next(ctx, req)

			attrs := []any{
				slog.String("request_id", requestID),
				slog.String("procedure", req.Spec().Procedure),
				slog.String("duration", time.Since(start).String()),
			}
			if err != nil {
				attrs = append(attrs, slog.String("error", err.Error()))
				log.ErrorContext(ctx, "rpc", attrs...)
			} else {
				log.InfoContext(ctx, "rpc", attrs...)
			}

			return resp, err
		}
	}
}
