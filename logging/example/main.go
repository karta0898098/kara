package main

import (
	"context"

	"github.com/karta0898098/kara/logging"
	"github.com/rs/zerolog/log"
)

func main() {
	logger := logging.Setup(logging.Config{
		Env:   "local",
		App:   "app",
		Debug: true,
	})

	logger.Debug().Msg("call")
	logger.Info().Msg("call")
	logger.Warn().Msg("call")
	logger.Error().Msg("call")
	logger.Trace().Msg("call")

	ctx := context.Background()
	ctx = WithValue(ctx, "trace_id", "test-trace-id")

	log.Ctx(ctx).Info().Msg("call")
}

func WithValue(ctx context.Context, key string, value interface{}) context.Context {
	ctx = context.WithValue(ctx, key, value)
	logger := log.With().Fields(map[string]interface{}{
		key: value,
	}).Logger()
	ctx = logger.WithContext(ctx)
	return ctx
}
