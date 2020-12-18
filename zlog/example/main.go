package main

import (
	"context"

	"github.com/karta0898098/kara/zlog"
	"github.com/rs/zerolog/log"
)

func main() {
	zlog.Setup(&zlog.Config{
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
	logger := log.With().Str("trace_id", "12345").Logger()
	ctx = logger.WithContext(ctx)

	log.Ctx(ctx).Info().Msg("call")
}
