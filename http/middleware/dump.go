package middleware

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func NewEchoDumpMiddleware() echo.MiddlewareFunc {
	return BodyDumpWithConfig(BodyDumpConfig{
		RequestSkipper: func(c echo.Context) bool {
			ctx := c.Request().Context()
			logger := log.Ctx(ctx)
			contentType := c.Request().Header.Get(echo.HeaderContentType)
			switch contentType {
			case "":
				return false
			case echo.MIMEApplicationJSON:
				return false
			case echo.MIMEApplicationJSONCharsetUTF8:
				return false
			default:
				logger.Info().Msgf("http request dump not support type %s", contentType)
				return true
			}
		},
		ResponseSkipper: func(c echo.Context) bool {
			ctx := c.Request().Context()
			logger := log.Ctx(ctx)
			contentType := c.Response().Header().Get(echo.HeaderContentType)
			switch contentType {
			case "":
				return false
			case echo.MIMEApplicationJSON:
				return false
			case echo.MIMEApplicationJSONCharsetUTF8:
				return false
			default:
				logger.Info().Msgf("http response dump not support type %s", contentType)
				return true
			}
		},
		RequestHandler: func(c echo.Context, req []byte) {
			ctx := c.Request().Context()
			logger := log.Ctx(ctx)
			if len(req) == 0 {
				logger.Info().
					Interface("data", nil).
					Msg("http request dump data.")
			} else {
				logger.Info().
					RawJSON("data", req).
					Msg("http request dump data.")
			}
		},
		ResponseHandler: func(c echo.Context, resp []byte) {
			ctx := c.Request().Context()
			logger := log.Ctx(ctx)

			if len(resp) == 0 {
				logger.Info().
					Interface("data", nil).
					Msg("http response dump data.")
				return
			} else {
				logger.Info().
					RawJSON("data", resp).
					Msg("http response dump data.")
			}
		},
	})
}

type (
	// BodyDumpConfig defines the config for BodyDump middleware.
	BodyDumpConfig struct {
		// RequestSkipper defines a function to skip middleware.
		RequestSkipper middleware.Skipper

		// ResponseSkipper defines a function to skip middleware.
		ResponseSkipper middleware.Skipper

		// Handler receives request and response payload.
		// Required.
		RequestHandler BodyDumpHandler

		ResponseHandler BodyDumpHandler
	}

	// BodyDumpHandler receives the request and response payload.
	BodyDumpHandler func(echo.Context, []byte)

	bodyDumpResponseWriter struct {
		io.Writer
		http.ResponseWriter
	}
)

var (
	// DefaultBodyDumpConfig is the default BodyDump middleware config.
	DefaultBodyDumpConfig = BodyDumpConfig{
		RequestSkipper: middleware.DefaultSkipper,
	}
)

// BodyDumpWithConfig returns a BodyDump middleware with config.
// See: `BodyDump()`.
func BodyDumpWithConfig(config BodyDumpConfig) echo.MiddlewareFunc {
	// Defaults
	if config.RequestHandler == nil {
		panic("echo: body-dump middleware requires a handler function")
	}

	if config.RequestHandler == nil {
		panic("echo: body-dump middleware requires a handler function")
	}

	if config.RequestSkipper == nil {
		config.RequestSkipper = DefaultBodyDumpConfig.RequestSkipper
	}

	if config.ResponseSkipper == nil {
		config.ResponseSkipper = DefaultBodyDumpConfig.RequestSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if !config.RequestSkipper(c) {
				// Request
				var reqBody []byte
				if c.Request().Body != nil { // Read
					reqBody, _ = ioutil.ReadAll(c.Request().Body)
				}
				c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset
				config.RequestHandler(c, reqBody)
			}

			err = next(c)

			// Response
			if !config.ResponseSkipper(c) {
				resBody := new(bytes.Buffer)
				mw := io.MultiWriter(c.Response().Writer, resBody)
				writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
				c.Response().Writer = writer
				c.Response().After(func() {
					config.ResponseHandler(c, resBody.Bytes())
				})
			}
			return
		}
	}
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
