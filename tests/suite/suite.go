package suite

import (
	"context"
	authv1 "github.com/TauAdam/sso/contracts/gen/go/sso"
	"github.com/TauAdam/sso/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Config     *config.Config
	AuthClient authv1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadConfigByPath("../config/tests.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	// Create a new gRPC client, use insecure credentials for testing
	conn, err := grpc.DialContext(context.Background(),
		net.JoinHostPort("localhost", strconv.Itoa(cfg.GRPC.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Config:     cfg,
		AuthClient: authv1.NewAuthClient(conn),
	}
}
