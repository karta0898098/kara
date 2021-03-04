package logging

import (
	"context"
	"net"

	"github.com/karta0898098/kara/tracer"
	"github.com/labstack/echo/v4"
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
		ctx = context.WithValue(ctx, tracer.TraceIDKey, traceID)

		logger :=
			log.With().
				Str("remote_addr", addr).
				Str("method", info.FullMethod).
				Str("trace_id", traceID).
				Logger()

		logger.Info().Msg("grpc access log.")

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

		traceID := ctx.Value(tracer.TraceIDKey).(string)

		md.Set(echo.HeaderXRequestID, traceID)

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
		meta := md.Get(echo.HeaderXRequestID)
		if len(meta) > 0 {
			requestID = meta[0]
		}
	}

	return requestID
}
