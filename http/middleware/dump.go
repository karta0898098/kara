package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func NewEchoDumpMiddleware() echo.MiddlewareFunc {
	return middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
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
