package grpc_auth

import (
	"context"
	"google.golang.org/grpc"
	goauth "pkg.tanyudii.me/go-pkg/go-auth"
)

func UnaryInterceptor(
	authService goauth.Service,
	tokenService goauth.TokenService,
	args ...ConfigFunc,
) grpc.UnaryServerInterceptor {
	svc := newService(authService, tokenService, args...)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx, err := svc.authenticate(ctx, info)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}
