package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// RequestIDFromContext 從 ctx 中取得 request id, 如果沒有即時產生一個
func RequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(echo.HeaderXRequestID).(string)
	if !ok {
		// 產生 requestID 並傳下去
		requestID = uuid.New().String()
		return requestID
	}
	return requestID
}

// NewRequestIDMiddleware Default returns the location middleware with default configuration.
func NewRequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = uuid.New().String()
			}
			c.Request().Header.Set(echo.HeaderXRequestID, requestID)

			logger := log.With().Str("request_id", requestID).Logger()
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