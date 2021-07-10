package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karta0898098/kara/logging"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewEchoDumpMiddleware(t *testing.T) {
	var (
		jsonResponse = func(c echo.Context) error {
			return c.JSON(http.StatusOK, echo.Map{
				"status": "ok",
			})
		}
		stringResponse = func(c echo.Context) error {
			return c.String(http.StatusOK, "not support")
		}
	)

	logger, tests := logging.SetupWithOption(
		logging.WithDebug(true),
		logging.WithLevel(logging.DebugLevel),
	), []struct {
		name        string
		contentType string
		request     echo.Map
		response    echo.HandlerFunc
	}{
		{
			name:        "Success",
			contentType: echo.MIMEApplicationJSON,
			request:     echo.Map{"id": 1},
			response:    jsonResponse,
		},
		{
			name:        "NotSupportRequest",
			contentType: echo.MIMEApplicationForm,
			request:     echo.Map{"id": 1},
			response:    jsonResponse,
		},
		{
			name:        "NoData",
			contentType: "",
			request:     nil,
			response:    jsonResponse,
		},
		{
			name:        "NotSupportResponse",
			contentType: echo.MIMEApplicationJSON,
			request:     echo.Map{"id": 1},
			response:    stringResponse,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := logger.WithContext(context.Background())

			e := echo.New()
			payload, _ := json.Marshal(&tt.request)
			req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(payload))
			resp := httptest.NewRecorder()

			req.Header.Set(echo.HeaderContentType, tt.contentType)
			c := e.NewContext(req, resp)
			c.SetRequest(c.Request().WithContext(ctx))

			h := NewEchoDumpMiddleware()(tt.response)
			assert.NoError(t, h(c))
		})
	}
}
