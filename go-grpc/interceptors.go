package go_grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gotex "pkg.tanyudii.me/go-pkg/go-tex"
	"runtime/debug"
	"strings"
	"time"
)

func RecoveryUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Unknown, "unexpected error happened")
				fmt.Printf("panic recovered: %v; stacktrace: %s\n", r, string(debug.Stack()))
			}
		}()
		return handler(ctx, req)
	}
}

func RequestIDUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		md := gotex.FromIncoming(ctx)
		if md.Get(strings.ToLower(gotex.RequestHeaderKeyRequestID)) == "" {
			requestID := fmt.Sprintf("%s-%d", uuid.NewString(), time.Now().Unix())
			md.Set(strings.ToLower(gotex.RequestHeaderKeyRequestID), requestID)
			ctx = gotex.NewContext(ctx, gotex.NewGotex(md))
		}
		return handler(ctx, req)
	}
}

func AcceptLangUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		md := gotex.FromIncoming(ctx)
		if acceptLang := md.Get(strings.ToLower("grpcgateway-" + gotex.RequestHeaderKeyAcceptLanguage)); acceptLang != "" {
			md.Set(strings.ToLower(gotex.RequestHeaderKeyAcceptLanguage), acceptLang)
			ctx = gotex.NewContext(ctx, gotex.NewGotex(md))
		}
		return handler(ctx, req)
	}
}
