package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func NewEcho(config *Config) *echo.Echo {

	e := echo.New()

	if config.Mode == "release" {
		e.Debug = true
		e.HideBanner = true
		e.HidePort = false

	} else {
		e.Debug = false
		e.HideBanner = true
		e.HidePort = true
	}

	return e
}

func RunEcho(engine *echo.Echo, config *Config,lifecycle fx.Lifecycle) *echo.Echo  {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			go func() {
				err = engine.Start(config.Port)
				if err != nil {
					log.Error().Msgf("Error echo server, err: %v", err)
				}
			}()
			return err
		},
		OnStop: func(ctx context.Context) error {
			return engine.Shutdown(ctx)
		},
	})
	return engine
}