package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"

	"github.com/ABCproger/card-validator/config"
	"github.com/ABCproger/card-validator/gen/card/v1/cardv1connect"
	controllerv1 "github.com/ABCproger/card-validator/internal/controller/connect/v1"
	uccard "github.com/ABCproger/card-validator/internal/usecase/card"
	"github.com/ABCproger/card-validator/pkg/httpserver"
	"github.com/ABCproger/card-validator/pkg/interceptor"
	"github.com/ABCproger/card-validator/pkg/logger"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	uc := uccard.New()
	handler := controllerv1.NewCardHandler(uc)

	mux := http.NewServeMux()

	path, connectHandler := cardv1connect.NewCardServiceHandler(
		handler,
		connect.WithInterceptors(interceptor.NewLoggingInterceptor(log)),
		connect.WithCompressMinBytes(1024),
	)
	mux.Handle(path, connectHandler)

	reflector := grpcreflect.NewStaticReflector(cardv1connect.CardServiceName)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	addr := fmt.Sprintf(":%s", cfg.HTTP.Port)
	srv := httpserver.New(mux, addr)

	log.Info("server started", slog.String("addr", addr), slog.String("env", cfg.App.Env))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-srv.Notify():
		log.Error("server error", slog.Any("err", err))
	case sig := <-quit:
		log.Info("shutting down", slog.String("signal", sig.String()))
	}

	if err := srv.Shutdown(); err != nil {
		log.Error("shutdown error", slog.Any("err", err))
		os.Exit(1)
	}
	log.Info("server stopped")
}
