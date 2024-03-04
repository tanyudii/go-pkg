package go_grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	ErrUnsupportedNetwork = fmt.Errorf("unsupported network")
)

const (
	defaultMaxCallRcvMsgSize  = 1024 * 1024 * 50 //50MB
	defaultMaxCallSendMsgSize = 1024 * 1024 * 50 //50MB
)

func ClientConn(addr string) grpc.ClientConnInterface {
	opts := getDialOpts()
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		panic(err)
	}
	return conn
}

func getDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
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

func dialTCP(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, addr, getDialOpts()...)
}
