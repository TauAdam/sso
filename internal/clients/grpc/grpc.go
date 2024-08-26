package grpc

import (
	"context"
	"fmt"
	authv1 "github.com/TauAdam/sso/contracts/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type Client struct {
	api authv1.AuthClient
	log *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, address string, timeout time.Duration, retriesNumber int) (*Client, error) {
	const op = "grpc.New"

	retryOptions := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesNumber)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s failed to dial: %w", op, err)
	}

	return &Client{
		api: authv1.NewAuthClient(conn),
	}, nil
}
