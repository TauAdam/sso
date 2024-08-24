package app

import (
	grpcapp "github.com/TauAdam/sso/internal/app/grpc"
	"log/slog"
	"time"
)

type Application struct {
	GRPCServer *grpcapp.Application
}

func New(grpcPort int, log *slog.Logger, tokenTTL time.Duration, storagePath string) *Application {
	//	TODO initialize the storage

	//	 TODO initialize the auth service

	grpcApp := grpcapp.New(grpcPort, log)

	return &Application{
		GRPCServer: grpcApp,
	}
}
