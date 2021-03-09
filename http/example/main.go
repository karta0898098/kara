package main

import (
	"github.com/karta0898098/kara/errors"
	"github.com/karta0898098/kara/http"
	"github.com/karta0898098/kara/zlog"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	zlog.Setup(zlog.Config{
		Env:   "local",
		AppID: "example",
		Debug: true,
	})

	e := http.NewEcho(http.Config{
		Mode: "debug",
		Port: ":8080",
		Dump: false,
	})

	e.POST("/api/v1/login", func(c echo.Context) error {
		return errors.ErrInvalidInput.Build("testing error")
	})

	log.Fatal().Err(e.Start(":8080")).Msg("start server failed")
}
