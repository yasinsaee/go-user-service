package middleware

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt"
	jwt2 "github.com/yasinsaee/go-user-service/pkg/jwt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if strings.Contains(info.FullMethod, "Login") {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		tokenList := md["authorization"]
		if len(tokenList) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is required")
		}
		tokenStr := strings.TrimPrefix(tokenList[0], "Bearer ")

		pubKey, err := jwt2.GetPublicKey()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to load public key: %v", err)
		}

		claims, err := jwt.ParseWithClaims(tokenStr, &jwt2.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return pubKey, nil
		})
		if err != nil || !claims.Valid {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		newCtx := context.WithValue(ctx, "user", claims)
		return handler(newCtx, req)
	}
}
