package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			traceID := c.Request().Header.Get(echo.HeaderXRequestID)

			start := time.Now()
			err := next(c)
			stop := time.Now()

			var logger *zerolog.Event

			status := c.Response().Status
			if status >= 500 {
				logger = log.Error()
			} else if status >= 400 {
				logger = log.Info()
			} else {
				logger = log.Info()
			}

			logger.
				Str("method", c.Request().Method).
				Str("uri", c.Request().RequestURI).
				Str("trace_id", traceID).
				Str("latency_human", stop.Sub(start).String()).
				Int("status", status).
				Err(err)

			logger.Msg("http access log.")
			return nil
		}
	}
}
