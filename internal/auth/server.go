package auth

import (
	"context"

	authpb "github.com/Maraei/calculator-on-go/api/api"
)

type AuthServer struct {
	authpb.UnimplementedAuthCalculatorServiceServer
	store *Store
}

func NewAuthServer(store *Store) *AuthServer {
	return &AuthServer{store: store}
}

func (s *AuthServer) Register(ctx context.Context, req *authpb.AuthRequest) (*authpb.AuthResponse, error) {
	err := s.store.CreateUser(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &authpb.AuthResponse{Message: "User registered successfully"}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.AuthRequest) (*authpb.TokenResponse, error) {
	err := s.store.ValidateUser(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	token, err := GenerateToken(req.Username)
	if err != nil {
		return nil, err
	}

	return &authpb.TokenResponse{Token: token}, nil
}

func (s *AuthServer) Validate(ctx context.Context, req *authpb.TokenRequest) (*authpb.ValidateResponse, error) {
	username, err := ValidateToken(req.Token)
	if err != nil {
		return &authpb.ValidateResponse{Valid: false}, nil
	}

	return &authpb.ValidateResponse{
		Valid:    true,
		Username: username,
	}, nil
}
