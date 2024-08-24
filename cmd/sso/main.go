package main

import (
	"github.com/TauAdam/sso/internal/app"
	"github.com/TauAdam/sso/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoadConfig()

	log := prepareLogger(cfg.Env)

	log.Info("starting application", slog.Any("config", cfg))

	application := app.New(cfg.GRPC.Port, log, cfg.TokenTTL, cfg.StoragePath)

	go application.GRPCServer.MustRun()

	//	 Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sign := <-quit

	log.Info("shutting down application", slog.String(
		"signal",
		sign.String(),
	))

	application.GRPCServer.Stop()

	log.Info("application stopped")
}

const (
	environmentLocal = "local"
	environmentDev   = "dev"
	environmentProd  = "prod"
)

func prepareLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case environmentLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case environmentDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case environmentProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
