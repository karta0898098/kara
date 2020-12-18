package logging

import (
	"context"
	"net"

	"github.com/google/uuid"
	"github.com/karta0898098/kara/metrics"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// UnaryServerInterceptor returns a new unary server logging and set request id.
func UnaryServerInterceptor(dump bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var (
			traceID string
			addr    string
		)

		if pr, ok := peer.FromContext(ctx); ok {
			if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
				addr = tcpAddr.IP.String()
			} else {
				addr = pr.Addr.String()
			}
		}

		traceID = getTraceID(ctx)
		ctx = context.WithValue(ctx, metrics.DefaultTraceID, traceID)

		logger :=
			log.With().
				Str("remote_addr", addr).
				Str("method", info.FullMethod).
				Str("trace_id", traceID).
				Logger()

		logger.Info().Msg("grpc access log")

		requestDump(ctx, info, dump, logger, req, nil)
		resp, err := handler(ctx, req)
		replayDump(ctx, info, dump, logger, resp, err)
		return resp, err
	}
}

// UnaryClientInterceptor ...
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}

		traceID := metrics.GetTraceID(ctx)
		if traceID == "" {
			traceID = uuid.New().String()
		}
		md.Set(metrics.DefaultTraceID, traceID)

		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func getTraceID(ctx context.Context) string {
	var (
		requestID string
	)

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		meta := md.Get(metrics.DefaultTraceID)
		if len(meta) > 0 {
			requestID = meta[0]
		}
	}

	return requestID
}
