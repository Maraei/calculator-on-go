package auth

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey string

const userIDKey contextKey = "user_id"

func AuthMiddleware(secret string) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, status.Error(codes.Unauthenticated, "missing metadata")
        }

        authHeader := md["authorization"]
        if len(authHeader) == 0 || !strings.HasPrefix(authHeader[0], "Bearer ") {
            return nil, status.Error(codes.Unauthenticated, "missing or malformed token")
        }

        tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil // используем переданный секрет
        })
        if err != nil || !token.Valid {
            return nil, status.Error(codes.Unauthenticated, "invalid token")
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return nil, status.Error(codes.Unauthenticated, "invalid claims")
        }

        idFloat, ok := claims["id"].(float64)
        if !ok {
            return nil, status.Error(codes.Unauthenticated, "invalid user ID type")
        }
        userID := int(idFloat)

        ctx = context.WithValue(ctx, userIDKey, userID)

        return handler(ctx, req)
    }
}
