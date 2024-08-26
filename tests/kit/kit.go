package kit

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

type Kit struct {
	*testing.T
	Config     *config.Config
	AuthClient authv1.AuthClient
}

func New(t *testing.T) (context.Context, *Kit) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadConfigByPath("../config/local.yml")

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

	return ctx, &Kit{
		T:          t,
		Config:     cfg,
		AuthClient: authv1.NewAuthClient(conn),
	}
}
