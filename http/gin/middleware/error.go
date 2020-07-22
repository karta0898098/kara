package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func RecordErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Errors.Last()

		if err != nil {
			logFields := map[string]interface{}{}

			logFields["request_method"] = c.Request.Method
			logFields["request_url"] = c.Request.URL.String()

			c.Next()
			status := c.Writer.Status()
			logFields["response_status"] = status
			// 根據狀態碼用不同等級來紀錄
			logger := log.Ctx(c.Request.Context()).With().Fields(logFields).Logger()
			if status >= http.StatusInternalServerError {
				logger.Error().Msgf("%+v", err)
			} else if status >= http.StatusBadRequest {
				logger.Debug().Msgf("%+v", err)
			} else {
				logger.Debug().Msgf("%+v", err)
			}
			return
		}
		c.Next()
	}
}
