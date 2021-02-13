package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/karta0898098/kara/tracer"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// NewTracerMiddleware Default returns the location middleware with default configuration.
func NewTracerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			traceID := c.Request().Header.Get(echo.HeaderXRequestID)
			if traceID == "" {
				traceID = uuid.New().String()
				c.Request().Header.Set(echo.HeaderXRequestID, traceID)
			}

			// set context trace id
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, tracer.TraceIDKey, traceID)

			// set logger trace id
			logger := log.With().Str(tracer.TraceIDKey.ToString(), traceID).Logger()
			ctx = logger.WithContext(ctx)

			c.SetRequest(c.Request().WithContext(ctx))
			// Set X-Request-Id header
			c.Response().Writer.Header().Set(echo.HeaderXRequestID, traceID)
			return next(c)
		}
	}
}
