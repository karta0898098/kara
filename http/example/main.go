package main

import (
	"github.com/bwmarrin/snowflake"
	"github.com/karta0898098/kara/http"
	"github.com/karta0898098/kara/zlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func main() {
	node, _ := snowflake.NewNode(1)
	httpConfig := &http.Config{
		Mode: "debug",
		Port: ":8080",
		Dump: true,
	}

	logConfig := &zlog.Config{
		Env:   "local",
		AppID: "echo_test",
		Level: 0,
		Debug: true,
	}

	var router *echo.Echo
	app := fx.New(
		fx.Supply(httpConfig),
		fx.Supply(logConfig),
		fx.Supply(node),
		fx.Provide(http.NewEcho),
		fx.Invoke(zlog.Setup),
		fx.Invoke(http.RunEcho),
		fx.Populate(&router),
	)

	// router.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
	// 	Skipper: middleware.DefaultSkipper,
	// 	Handler: func(ctx echo.Context, req []byte, resp []byte) {
	// 		var (
	// 			reqDict  *zerolog.Event
	// 			respDict *zerolog.Event
	// 		)
	//
	// 		reqDict = zerolog.Dict()
	// 		respDict = zerolog.Dict()
	//
	// 		if len(req) > 0 {
	// 			reqDict.RawJSON("body", req)
	// 		}
	// 		log.Ctx(ctx.Request().Context()).
	// 			Info().
	// 			Dict("dump", reqDict).
	// 			Msg("http request dump data.")
	//
	//
	// 		if len(resp) > 0 {
	// 			respDict.RawJSON("body", resp)
	// 		}
	// 		log.Ctx(ctx.Request().Context()).
	// 			Info().
	// 			Dict("dump", respDict).
	// 			Msg("http response dump data.")
	// 	},
	// }))

	router.POST("/ping", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"hello": "world",
		})
	})
	app.Run()
}
