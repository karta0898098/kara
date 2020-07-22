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
		requestID := c.GetHeader(headerXRequestID)

		if requestID == "" {
			requestID = uuid.New().String()
			c.Header(headerXRequestID, requestID)
		}

		logger := log.With().Str("request_id", requestID).Logger()
		ctx := logger.WithContext(c.Request.Context())
		ctx = context.WithValue(ctx, headerXRequestID, requestID)
		ctx = logger.WithContext(ctx)
		// Set X-Request-Id header
		c.Header(headerXRequestID, requestID)
		c.Next()
	}
}
