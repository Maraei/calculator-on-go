package auth

import (
	"context"
	"fmt"
	"database/sql"

	authpb "github.com/Maraei/calculator-on-go/api/api"
	_ "github.com/mattn/go-sqlite3"
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

	userID, err := s.store.GetUserIDByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	token, err := GenerateToken(userID)
	if err != nil {
		return nil, err
	}

	return &authpb.TokenResponse{Token: token}, nil
}


func (s *Store) GetUserIDByUsername(username string) (int, error) {
	var userID int
	err := s.db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user not found")
		}
		return 0, err
	}
	return userID, nil
}


func (s *AuthServer) Validate(ctx context.Context, req *authpb.TokenRequest) (*authpb.ValidateResponse, error) {
    userID, err := ValidateToken(req.Token) 
    if err != nil {
        return &authpb.ValidateResponse{Valid: false}, nil
    }

    username, err := s.store.GetUsernameByID(userID)
    if err != nil {
        return &authpb.ValidateResponse{Valid: false}, nil
    }

    return &authpb.ValidateResponse{
        Valid:    true,
        Username: username,
    }, nil
}

func (s *Store) GetUsernameByID(userID int) (string, error) {
    var username string
    err := s.db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("user not found")
        }
        return "", err
    }
    return username, nil
}
