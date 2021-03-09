package main

import (
	"context"

	"github.com/karta0898098/kara/zlog"
	"github.com/rs/zerolog/log"
)

func main() {
	zlog.Setup(zlog.Config{
		Env:   "local",
		AppID: "app",
		Level: -1,
		Debug: true,
	})

	log.Debug().Msg("call")
	log.Info().Msg("call")
	log.Warn().Msg("call")
	log.Error().Msg("call")
	log.Trace().Msg("call")

	ctx := context.Background()
	ctx = WithValue(ctx, "trace_id", "1234567")
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
