package grpc

import (
	"time"

	"github.com/karta0898098/kara/grpc/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// NewClient new grpc client
func NewClient(host string) (*grpc.ClientConn, error) {
	var (
		interceptors []grpc.UnaryClientInterceptor
	)

	interceptors = []grpc.UnaryClientInterceptor{
		logging.UnaryClientInterceptor(),
	}

	options := grpc.WithChainUnaryInterceptor(interceptors...)

	conn, err := grpc.Dial(host, grpc.WithInsecure(),
		grpc.WithInitialConnWindowSize(256*1024),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second,
			Timeout:             18 * time.Second,
			PermitWithoutStream: true,
		}),
		options,
	)
	return conn, err
}
