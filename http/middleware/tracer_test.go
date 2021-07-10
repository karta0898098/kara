package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/karta0898098/kara/logging"
	"github.com/karta0898098/kara/tracer"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestNewTracerMiddleware(t *testing.T) {
	// setup default logger
	logging.Setup(logging.Config{
		Env:   "local",
		App:   "test",
		Debug: true,
	})

	e := echo.New()

	tests := []struct {
		name    string
		traceID string
	}{
		{
			name:    "SuccessEqTraceID",
			traceID: "test_trace_id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()
			req.Header.Set(echo.HeaderXRequestID, tt.traceID)
			c := e.NewContext(req, resp)
			h := NewTracerMiddleware()(func(c echo.Context) error {
				ctx := c.Request().Context()
				assert.Equal(t, tt.traceID, ctx.Value(tracer.TraceIDKey))
				// check log trace id eq context.Context trace id
				log.Ctx(ctx).Info().Msgf("log trace id = %v", ctx.Value(tracer.TraceIDKey))
				return nil
			})

			assert.NoError(t, h(c))
			assert.Equal(t, tt.traceID, resp.Header().Get(echo.HeaderXRequestID))
		})
	}
}

func TestNewTracerMiddlewareWithConcurrency(t *testing.T) {
	// setup default logger
	logging.Setup(logging.Config{
		Env:   "local",
		App:   "test",
		Debug: true,
	})

	e := echo.New()

	// simulation 10 request concurrency
	task := 10
	wg := new(sync.WaitGroup)
	wg.Add(task)

	for i := 0; i < task; i++ {
		go func(i int) {
			expTraceID := "trace_id_" + strconv.Itoa(i)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()
			req.Header.Set(echo.HeaderXRequestID, expTraceID)

			c := e.NewContext(req, resp)

			h := NewTracerMiddleware()(func(c echo.Context) error {
				ctx := c.Request().Context()
				assert.Equal(t, expTraceID, ctx.Value(tracer.TraceIDKey))

				// check log trace id eq context.Context trace id
				log.Ctx(ctx).Info().Msgf("log trace id = %v", ctx.Value(tracer.TraceIDKey))
				return nil
			})
			assert.NoError(t, h(c))
			assert.Equal(t, expTraceID, resp.Header().Get(echo.HeaderXRequestID))
			wg.Done()
		}(i)
	}
	wg.Wait()
}
