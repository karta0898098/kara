package http

import (
	"context"
	"net/http"

	"github.com/karta0898098/kara/errors"
	"github.com/karta0898098/kara/http/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

// NewEcho new echo http engine constructor
func NewEcho(config Config) *echo.Echo {
	echo.NotFoundHandler = EchoNotFoundHandler

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
	e.Pre(middleware.NewTracerMiddleware())
	e.Use(middleware.NewLoggerMiddleware())
	e.Use(middleware.NewErrorHandlingMiddleware())
	e.Use(middleware.RecordErrorMiddleware())

	if config.Dump {
		e.Use(middleware.NewEchoDumpMiddleware())
	}

	return e
}

// RunEcho for use uber fx to start http service
func RunEcho(engine *echo.Echo, config Config, lifecycle fx.Lifecycle) *echo.Echo {
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

	appException := errors.TryConvert(err)
	if appException == nil {
		status, payload := errors.ErrInternal.ToViewModel()
		_ = c.JSON(status, payload)
		return
	}

	status, payload := appException.ToViewModel()

	_ = c.JSON(status, payload)
}

// EchoNotFoundHandler responds not found response.
func EchoNotFoundHandler(c echo.Context) error {
	status, payload := errors.ErrPageNotFound.ToViewModel()
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
