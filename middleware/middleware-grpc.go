package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const ContextUserIDKey contextKey = "id"

func JWTUnaryInterceptor(jwtSign string, publicMethods ...string) grpc.UnaryServerInterceptor {
	public := toSet(publicMethods)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if isPublic(info.FullMethod, public) {
			return handler(ctx, req)
		}
		newCtx, err := authorizeGRPC(ctx, jwtSign)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func JWTStreamInterceptor(jwtSign string, publicMethods ...string) grpc.StreamServerInterceptor {
	public := toSet(publicMethods)
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if isPublic(info.FullMethod, public) {
			return handler(srv, ss)
		}
		if _, err := authorizeGRPC(ss.Context(), jwtSign); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func toSet(methods []string) map[string]struct{} {
	set := make(map[string]struct{}, len(methods))
	for _, m := range methods {
		set[m] = struct{}{}
	}
	return set
}

func isPublic(fullMethod string, public map[string]struct{}) bool {
	if strings.HasPrefix(fullMethod, "/grpc.reflection.") {
		return true
	}
	_, ok := public[fullMethod]
	return ok
}

func authorizeGRPC(ctx context.Context, jwtSign string) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "missing metadata")
	}
	vals := md.Get("authorization")
	if len(vals) == 0 {
		return ctx, status.Error(codes.Unauthenticated, "missing authorization token")
	}

	signature := strings.Split(vals[0], " ")
	if len(signature) < 2 || signature[0] != "Bearer" {
		return ctx, status.Error(codes.Unauthenticated, "authorization must be a Bearer token")
	}

	claim := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(signature[1], claim, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSign), nil
	})
	if err != nil {
		return ctx, status.Error(codes.Unauthenticated, "invalid token")
	}

	if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || method != jwt.SigningMethodHS256 {
		return ctx, status.Error(codes.Unauthenticated, "invalid signing method")
	}

	expAt, err := claim.GetExpirationTime()
	if err != nil || expAt == nil || time.Now().After(expAt.Time) {
		return ctx, status.Error(codes.Unauthenticated, "token expired")
	}

	if userIDFloat, ok := claim["id"].(float64); ok {
		ctx = context.WithValue(ctx, ContextUserIDKey, int(userIDFloat))
	}
	return ctx, nil
}
