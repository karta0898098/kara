package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karta0898098/kara/zlog"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewEchoDumpMiddleware(t *testing.T) {
	zlog.Setup(&zlog.Config{
		Env:   "local",
		AppID: "test",
		Debug: true,
	})

	tests := []struct {
		name        string
		contentType string
		request     echo.Map
	}{
		{
			name:        "success",
			contentType: echo.MIMEApplicationJSON,
			request:     echo.Map{"id": 1},
		},
		{
			name:        "nonsupport",
			contentType: echo.MIMEApplicationForm,
			request:     echo.Map{"id": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			payload, _ := json.Marshal(&tt.request)
			req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(payload))
			resp := httptest.NewRecorder()

			req.Header.Set(echo.HeaderContentType, tt.contentType)

			c := e.NewContext(req, resp)
			h := NewEchoDumpMiddleware()(func(c echo.Context) error {
				return c.JSON(http.StatusOK, echo.Map{
					"status": "ok",
				})
			})
			assert.NoError(t, h(c))
		})
	}
}
