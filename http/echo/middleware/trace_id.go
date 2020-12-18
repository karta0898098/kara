package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// NewRequestIDMiddleware Default returns the location middleware with default configuration.
func NewRequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = uuid.New().String()
			}
			c.Request().Header.Set(echo.HeaderXRequestID, requestID)

			logger := log.With().Str("trace_id", requestID).Logger()
			ctx := logger.WithContext(c.Request().Context())
			ctx = context.WithValue(ctx, echo.HeaderXRequestID, requestID)
			c.SetRequest(c.Request().WithContext(ctx))
			ctx = logger.WithContext(ctx)
			// Set X-Request-Id header
			c.Response().Writer.Header().Set(echo.HeaderXRequestID, requestID)
			return next(c)
		}
	}
}