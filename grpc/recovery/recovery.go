package recovery

import (
	"context"
	"fmt"
	"runtime"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a new unary server recovery for panic recovery.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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
