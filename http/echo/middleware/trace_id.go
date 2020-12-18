package middleware

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// NewTraceMiddleware Default returns the location middleware with default configuration.
func NewTraceMiddleware(node *snowflake.Node) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			traceID := c.Request().Header.Get(echo.HeaderXRequestID)
			if traceID == "" {
				traceID = node.Generate().Base58()
				c.Request().Header.Set(echo.HeaderXRequestID, traceID)
			}

			logger := log.With().Str("trace_id", traceID).Logger()
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, echo.HeaderXRequestID, traceID)
			ctx = context.WithValue(ctx, "trace_id", traceID)
			ctx = logger.WithContext(ctx)

			c.SetRequest(c.Request().WithContext(ctx))
			// Set X-Request-Id header
			c.Response().Writer.Header().Set(echo.HeaderXRequestID, traceID)
			return next(c)
		}
	}
}
