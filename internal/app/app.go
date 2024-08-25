package app

import (
	grpcapp "github.com/TauAdam/sso/internal/app/grpc"
	"github.com/TauAdam/sso/internal/services/auth"
	"github.com/TauAdam/sso/internal/services/storage/sqlite"
	"log/slog"
	"time"
)

type Application struct {
	GRPCServer *grpcapp.Application
}

// Application puts together all the components

func New(grpcPort int, log *slog.Logger, tokenTTL time.Duration, storagePath string) *Application {
	storageService, err := sqlite.New(storagePath)
	if err != nil {
		panic("failed to create storage service" + err.Error())
	}

	authService := auth.New(log, tokenTTL, storageService, storageService, storageService)

	grpcApp := grpcapp.New(grpcPort, log, authService)

	return &Application{
		GRPCServer: grpcApp,
	}
}
