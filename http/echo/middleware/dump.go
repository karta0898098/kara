package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewEchoDumpMiddleware() echo.MiddlewareFunc  {
	return middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Skipper: middleware.DefaultSkipper,
		Handler: func(ctx echo.Context, req []byte, resp []byte) {
			var (
				reqDict  *zerolog.Event
				respDict *zerolog.Event
			)

			reqDict = zerolog.Dict()
			respDict = zerolog.Dict()

			if len(req) > 0 {
				reqDict.RawJSON("body", req)
			}
			log.Ctx(ctx.Request().Context()).
				Info().
				Dict("dump", reqDict).
				Msg("http request dump data.")


			if len(resp) > 0 {
				respDict.RawJSON("body", resp)
			}
			log.Ctx(ctx.Request().Context()).
				Info().
				Dict("dump", respDict).
				Msg("http response dump data.")
		},
	})
}
