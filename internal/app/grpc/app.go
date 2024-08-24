package grpcapp

import (
	"fmt"
	authgprc "github.com/TauAdam/sso/internal/grpc/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type Application struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(port int, log *slog.Logger) *Application {

	gRPCServer := grpc.NewServer()

	authgprc.RegisterServer(gRPCServer)

	return &Application{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *Application) Start() error {
	const op = "grpcapp.Application.Start"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("failed to listen %s: %w", op, err)
	}

	log.Info("grpc server started", slog.String("address", fmt.Sprintf(":%d", listener.Addr().String())))

	if err := a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve %s: %w", op, err)
	}

	return nil
}

func (a *Application) Stop() {
	const op = "grpcapp.Application.Stop"

	a.log.With(slog.String("op", op)).Info("stopping grpc server on port:", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}

func (a *Application) MustRun() {
	if err := a.Start(); err != nil {
		panic("failed to start grpc server: " + err.Error())
	}
}
