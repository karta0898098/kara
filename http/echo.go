package http

import (
	"context"
	"net/http"

	"github.com/bwmarrin/snowflake"
	"github.com/go-playground/validator/v10"
	"github.com/karta0898098/kara/exception"
	"github.com/karta0898098/kara/http/echo/middleware"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"go.uber.org/fx"
)

// NewEcho new echo http engine constructor
func NewEcho(config *Config, node *snowflake.Node) *echo.Echo {
	echo.NotFoundHandler = EchoNotFoundHandler
	echo.MethodNotAllowedHandler = EchoNotFoundHandler

	e := echo.New()
	e.Validator = NewEchoValidator()

	if config.Mode == "release" {
		e.Debug = false
		e.HideBanner = true
		e.HidePort = true

	} else {
		e.Debug = true
		e.HideBanner = false
		e.HidePort = false
	}

	e.HTTPErrorHandler = EchoErrorHandler
	e.Pre(middleware.NewTraceMiddleware(node))
	e.Use(middleware.NewLoggerMiddleware())
	e.Use(middleware.NewErrorHandlingMiddleware())
	e.Use(middleware.NewCORS())
	e.Use(middleware.RecordErrorMiddleware)

	if config.Dump {
		e.Use(middleware.NewEchoDumpMiddleware())
	}

	return e
}

// RunEcho for use uber fx to start http service
func RunEcho(engine *echo.Echo, config *Config, lifecycle fx.Lifecycle) *echo.Echo {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			go func() {
				err = engine.Start(config.Port)
				if err != nil {
					if errors.Is(err, http.ErrServerClosed) {
						log.Info().Msg("http server close.")
						return
					}
					log.Error().Err(err).Msg("run echo http server failed.")
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

// EchoErrorHandler for handle error to http error
func EchoErrorHandler(err error, c echo.Context) {
	if err == nil {
		return
	}

	echoError, ok := err.(*echo.HTTPError)
	if ok {
		_ = c.JSON(echoError.Code, echoError)
		return
	}

	appException := exception.TryConvert(err)
	if appException == nil {
		_ = c.JSON(http.StatusInternalServerError, exception.ErrInternal)
		return
	}

	status, payload := appException.ToView()

	_ = c.JSON(status, payload)
}

// EchoNotFoundHandler responds not found response.
func EchoNotFoundHandler(c echo.Context) error {
	status, payload := exception.ErrPageNotFound.ToView()
	return c.JSON(status, payload)
}

// EchoValidator fot echo default validator
type EchoValidator struct {
	validator *validator.Validate
}

// NewEchoValidator new echo validator
func NewEchoValidator() *EchoValidator {
	return &EchoValidator{validator: validator.New()}
}

// Validate for echo validator interface
func (e *EchoValidator) Validate(i interface{}) error {
	return e.validator.Struct(i)
}
