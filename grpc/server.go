package grpc

import (
	"context"
	"net"
	"time"

	"github.com/karta0898098/kara/grpc/logging"
	"github.com/karta0898098/kara/grpc/recovery"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func newTCPServer(config *Config) (net.Listener, error) {
	listen, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Error().Msgf("grpc bind failed %v", err)
	}
	return listen, err
}

// NewGRPC new grpc server and add default interceptor
func NewGRPC(config *Config) (*grpc.Server, net.Listener, error) {
	var (
		interceptors []grpc.UnaryServerInterceptor
		listener     net.Listener
		err          error
	)
	listener, err = newTCPServer(config)
	if err != nil {
		return nil, nil, err
	}

	interceptors = []grpc.UnaryServerInterceptor{
		recovery.UnaryServerInterceptor(),
		logging.UnaryServerInterceptor(config.RequestDump),
	}

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
	return server, listener, err
}

// RunGRPC start grpc server by use uber fx
func RunGRPC(listener net.Listener, service *grpc.Server, lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Info().Msgf("starting grpc service listen on %s", listener.Addr().String())
				if err := service.Serve(listener); err != nil {
					log.Error().Msgf("failed to start grpc service: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msgf("stopping grpc service.")
			service.GracefulStop()
			return listener.Close()
		},
	})
}
