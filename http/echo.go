package http

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/karta0898098/kara/exception"
	"github.com/karta0898098/kara/http/echo/middleware"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"net/http"
)

func NewEcho(config *Config) *echo.Echo {
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
	e.Pre(middleware.NewRequestIDMiddleware())
	e.Use(middleware.NewLoggerMiddleware())
	e.Use(middleware.NewErrorHandlingMiddleware())
	e.Use(middleware.NewCORS())
	e.Use(middleware.RecordErrorMiddleware)

	return e
}

func RunEcho(engine *echo.Echo, config *Config, lifecycle fx.Lifecycle) *echo.Echo {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			go func() {
				err = engine.Start(config.Port)
				if err != nil {
					log.Error().Err(err).Msg("start echo server failed.")
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

func EchoErrorHandler(err error, c echo.Context) {
	if err == nil {
		return
	}

	echoError, ok := err.(*echo.HTTPError)
	if ok {
		_ = c.JSON(echoError.Code, echoError)
		return
	}

	causeError := errors.Cause(err)
	appError, ok := causeError.(*exception.Exception)
	if !ok || appError == nil {
		_ = c.JSON(http.StatusInternalServerError, exception.ErrServerInternal)
		return
	}

	_ = c.JSON(appError.Status, map[string]interface{}{
		"code":    appError.Code,
		"message": appError.Message,
		"details": appError.Details,
	})
}

// EchoNotFoundHandler responds not found response.
func EchoNotFoundHandler(c echo.Context) error {
	return c.JSON(http.StatusNotFound, exception.ErrPageNotFound)
}

type EchoValidator struct {
	validator *validator.Validate
}

func NewEchoValidator() *EchoValidator {
	return &EchoValidator{validator: validator.New()}
}

func (e *EchoValidator) Validate(i interface{}) error {
	return e.validator.Struct(i)
}
