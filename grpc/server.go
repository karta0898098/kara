package grpc

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"runtime"
	"time"
)

var Module = fx.Provide(
	NewTCPServer,
	NewGRPC,
)

func NewTCPServer(config *Config) (net.Listener, error) {
	listen, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Error().Msgf("grpc bind failed %v", err)
	}
	return listen, err
}

func NewGRPC(config *Config) *grpc.Server {

	var interceptors []grpc.UnaryServerInterceptor

	interceptors = append(interceptors, unaryServerInterceptor())

	options := grpc.ChainUnaryInterceptor(interceptors...)

	server := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:              time.Duration(5) * time.Second, // Ping the client if it is idle for 5 seconds to ensure the connection is still active
				Timeout:           time.Duration(5) * time.Second, // Wait 5 second for the ping ack before assuming the connection is dead
				MaxConnectionIdle: 5 * time.Minute,
			},
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             time.Duration(2) * time.Second, // If a client pings more than once every 2 seconds, terminate the connection
				PermitWithoutStream: true,                           // Allow pings even when there are no active streams
			},
		),
		options,
	)
	return server
}

func RunGRPC(listener net.Listener, server *grpc.Server, lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			go func() {
				log.Info().Msgf("starting grpc server listen on %s", listener.Addr().String())
				if err := server.Serve(listener); err != nil {
					log.Error().Msgf("failed to start grpc server: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msgf("stopping grpc server.")
			server.GracefulStop()
			return listener.Close()
		},
	})
}

// UnaryServerInterceptor returns a new unary server interceptor for panic recovery.
func unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				var msg string
				for i := 2; ; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					msg = msg + fmt.Sprintf("%s:%d\n", file, line)
				}
				log.Error().Msgf("%s\n↧↧↧↧↧↧ PANIC ↧↧↧↧↧↧\n%s↥↥↥↥↥↥ PANIC ↥↥↥↥↥↥", r, msg)
			}
		}()
		return handler(ctx, req)
	}
}


// NewClient ... new grpc client
func NewClient(host string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure(),
		grpc.WithInitialConnWindowSize(256*1024),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second,
			Timeout:             18 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	return conn, err
}
