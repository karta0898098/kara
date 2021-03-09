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
		Skipper: func(c echo.Context) bool {
			contentType := c.Request().Header.Get(echo.HeaderContentType)
			switch contentType {
			case echo.MIMEApplicationJSON:
				return false
			case echo.MIMEApplicationJSONCharsetUTF8:
				return false
			default:
				log.Info().Msgf("http request dump not support type %s", contentType)
				return true
			}
		},
		Handler: func(c echo.Context, req []byte, resp []byte) {
			log.Info().
				RawJSON("body", req).
				Msg("http request dump data.")


			log.Info().
				RawJSON("body", resp).
				Msg("http response dump data.")
		},
	})
}

type (
	// BodyDumpConfig defines the config for BodyDump middleware.
	BodyDumpConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Handler receives request and response payload.
		// Required.
		Handler BodyDumpHandler
	}

	// BodyDumpHandler receives the request and response payload.
	BodyDumpHandler func(echo.Context, []byte, []byte)

	bodyDumpResponseWriter struct {
		io.Writer
		http.ResponseWriter
	}
)

var (
	// DefaultBodyDumpConfig is the default BodyDump middleware config.
	DefaultBodyDumpConfig = BodyDumpConfig{
		Skipper: middleware.DefaultSkipper,
	}
)

// BodyDumpWithConfig returns a BodyDump middleware with config.
// See: `BodyDump()`.
func BodyDumpWithConfig(config BodyDumpConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Handler == nil {
		panic("echo: body-dump middleware requires a handler function")
	}
	if config.Skipper == nil {
		config.Skipper = DefaultBodyDumpConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			// Request
			reqBody := []byte{}
			if c.Request().Body != nil { // Read
				reqBody, _ = ioutil.ReadAll(c.Request().Body)
			}
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

			// Response
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer
			c.Response().After(func() {
				// Callback
				config.Handler(c, reqBody, resBody.Bytes())
			})

			return next(c)
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
