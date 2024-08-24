package auth

import (
	"context"
	authv1 "github.com/TauAdam/sso/contracts/gen/go/sso"
	"google.golang.org/grpc"
)

type Server struct {
	authv1.UnimplementedAuthServer
}

func RegisterServer(s *grpc.Server) {
	authv1.RegisterAuthServer(s, &Server{})
}

func (s *Server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	panic("not implemented")
}

func (s *Server) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	panic("not implemented")
}

func (s *Server) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	panic("not implemented")
}
