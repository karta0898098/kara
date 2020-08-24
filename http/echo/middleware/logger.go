package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

func NewLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			requestID := c.Request().Header.Get(echo.HeaderXRequestID)
			start := time.Now()
			err := next(c)
			stop := time.Now()

			var logger *zerolog.Event

			status := c.Response().Status
			if status >= 500 {
				logger = log.Error()
			} else if status >= 400 {
				logger = log.Warn()
			} else {
				logger = log.Info()
			}

			logger.Str("method", c.Request().Method).
				Str("uri", c.Request().RequestURI).
				Str("request_id", requestID).
				Str("latency_human",stop.Sub(start).String()).
				Int("status", status).
				Err(err)

			logger.Msg("access log")
			return nil
		}
	}
}
