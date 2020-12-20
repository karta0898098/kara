package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const headerXRequestID = "X-Request-ID"

func NewRequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(headerXRequestID)

		if traceID == "" {
			traceID = uuid.New().String()
			c.Header(headerXRequestID, traceID)
		}

		logger := log.With().Str("trace_id", traceID).Logger()
		ctx := logger.WithContext(c.Request.Context())
		ctx = context.WithValue(ctx, headerXRequestID, traceID)
		ctx = logger.WithContext(ctx)
		// Set X-Request-Id header
		c.Header(headerXRequestID, traceID)
		c.Next()
	}
}
