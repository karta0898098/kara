package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/karta0898098/kara/errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// NewErrorHandlingMiddleware handles panic error
func NewErrorHandlingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					trace := make([]byte, 4096)
					runtime.Stack(trace, true)
					traceID := c.Request().Header.Get(echo.HeaderXRequestID)
					customFields := map[string]interface{}{
						"url":         c.Request().RequestURI,
						"stack_error": string(trace),
						"request_id":  traceID,
					}
					err, ok := r.(error)
					if !ok {
						if err == nil {
							err = fmt.Errorf("%v", r)
						} else {
							err = fmt.Errorf("%v", err)
						}
					}
					logger := log.With().Fields(customFields).Logger()
					logger.Error().Msgf("http: unknown error: %v", err)

					status, payload := errors.ErrInternal.ToRestfulView()
					_ = c.JSON(status, payload)
				}
			}()
			return next(c)
		}
	}
}

// RecordErrorMiddleware provide error middleware
func RecordErrorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			logFields := map[string]interface{}{}

			// 紀錄 Request 資料
			req := c.Request()
			{
				logFields["method"] = req.Method
				logFields["uri"] = req.RequestURI
			}
			ctx := req.Context()

			// 紀錄 Response 資料
			resp := c.Response()
			resp.After(func() {
				logFields["status"] = resp.Status
				// 根據狀態碼用不同等級來紀錄
				logger := log.Ctx(ctx).With().Fields(logFields).Logger()
				if resp.Status >= http.StatusInternalServerError {
					logger.Error().Msgf("%+v", err)
				} else if resp.Status >= http.StatusBadRequest {
					logger.Debug().Msgf("%+v", err)
				} else {
					logger.Debug().Msgf("%+v", err)
				}
			})
		}
		return err
	}
}
