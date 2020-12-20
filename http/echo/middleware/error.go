package middleware

import (
	"fmt"
	"github.com/karta0898098/kara/exception"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"runtime"
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

					_ = c.JSON(500, errors.Wrap(exception.ErrInternal, err.Error()))
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
				logFields["request_method"] = req.Method
				logFields["request_url"] = req.URL.String()
			}
			ctx := req.Context()

			// 紀錄 Response 資料
			resp := c.Response()
			resp.After(func() {
				logFields["response_status"] = resp.Status
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
