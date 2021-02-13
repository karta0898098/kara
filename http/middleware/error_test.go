package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karta0898098/kara/errors"
	"github.com/karta0898098/kara/zlog"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func TestRecordErrorMiddleware(t *testing.T) {
	zlog.New(&zlog.Config{
		Env:   "local",
		AppID: "test",
		Debug: true,
	})

	tests := []struct {
		name string
		err  *errors.Exception
	}{
		{
			name: "internal",
			err:  errors.ErrInternal,
		},
		{
			name: "resourceNotFound",
			err:  errors.ErrResourceNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()

			e := echo.New()
			e.HTTPErrorHandler = func(err error, c echo.Context) {
				if err == nil {
					return
				}

				echoError, ok := err.(*echo.HTTPError)
				if ok {
					_ = c.JSON(echoError.Code, echoError)
					return
				}

				appException := errors.TryConvert(err)
				if appException == nil {
					status, payload := errors.ErrInternal.ToViewModel()
					_ = c.JSON(status, payload)
					return
				}

				status, payload := appException.ToViewModel()

				_ = c.JSON(status, payload)
			}
			e.Use(RecordErrorMiddleware())
			e.GET("/", func(c echo.Context) error {
				ctx := c.Request().Context()
				logger := log.With().Logger()
				ctx = logger.WithContext(ctx)
				c.SetRequest(c.Request().WithContext(ctx))
				return tt.err.Build(fmt.Errorf("test error"))
			})
			e.ServeHTTP(resp, req)
		})
	}
}
