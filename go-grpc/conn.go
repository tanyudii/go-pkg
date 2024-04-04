package go_grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrUnsupportedNetwork = fmt.Errorf("[ERROR]: Unsupported network")
)

const (
	defaultMaxCallRcvMsgSize  = 1024 * 1024 * 50 //50MB
	defaultMaxCallSendMsgSize = 1024 * 1024 * 50 //50MB
)

func ClientConn(addr string, s ...bool) grpc.ClientConnInterface {
	opts := getDialOpts(isSecure(s...))
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		panic(err)
	}
	return conn
}

func getDialOpts(s bool) []grpc.DialOption {
	creds := insecure.NewCredentials()
	if s {
		creds = credentials.NewTLS(&tls.Config{})
	}
	return []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(defaultMaxCallRcvMsgSize),
			grpc.MaxCallSendMsgSize(defaultMaxCallSendMsgSize),
		),
	}
}

func dial(network, addr string) (*grpc.ClientConn, error) {
	switch network {
	case "tcp":
		return dialTCP(context.Background(), addr)
	default:
		return nil, fmt.Errorf("%w: %v", ErrUnsupportedNetwork, network)
	}
}

func dialTCP(ctx context.Context, addr string, s ...bool) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, addr, getDialOpts(isSecure(s...))...)
}

func isSecure(s ...bool) bool {
	if len(s) > 0 {
		return s[0]
	}
	return false
}
