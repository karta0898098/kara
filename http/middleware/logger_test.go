package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karta0898098/kara/zlog"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewLoggerMiddleware(t *testing.T) {
	zlog.Setup(zlog.Config{
		Env:   "local",
		AppID: "test",
		Debug: true,
	})

	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "StatusOK",
			status: http.StatusOK,
		},
		{
			name:   "StatusInternal",
			status: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()
			c := e.NewContext(req, resp)

			h := NewLoggerMiddleware()(func(c echo.Context) error {
				return c.NoContent(tt.status)
			})
			assert.NoError(t, h(c))
			assert.Equal(t, tt.status, resp.Code)
		})
	}
}
