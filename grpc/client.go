package grpc

import (
	"github.com/karta0898098/kara/grpc/logging"

	"google.golang.org/grpc"
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

	conn, err := grpc.Dial(host, grpc.WithInsecure(), options)
	return conn, err
}
