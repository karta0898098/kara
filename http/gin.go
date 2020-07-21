package http

import (
	"context"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"net/http"
)



func NewGin(config *Config) *gin.Engine {

	gin.SetMode(config.Mode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(logger.SetLogger())

	return engine
}

func RunGin(engine *gin.Engine, config *Config, lifecycle fx.Lifecycle) *gin.Engine {

	//engine := NewGin(config)
	srv := &http.Server{
		Addr:    config.Port,
		Handler: engine,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			go func() {
				err = srv.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					log.Info().Msgf("shutting down server %v", err)
				}
			}()
			return err
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return engine
}


