package main

import (
	"github.com/karta0898098/kara/errors"
	"github.com/karta0898098/kara/http"
	"github.com/karta0898098/kara/logging"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	logging.Setup(logging.Config{
		Env:   "local",
		App:   "example",
		Debug: true,
	})

	e := http.NewEcho(http.Config{
		Mode: "debug",
		Port: ":8080",
		Dump: true,
	})

	e.GET("/api/v1/login", func(c echo.Context) error {
		return errors.ErrInvalidInput.Build("testing error")
	})

	e.GET("/api/v1/test", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	e.GET("/api/v1/panic", func(c echo.Context) error {
		panic("call")
	})

	log.Fatal().Err(e.Start(":8080")).Msg("start server failed")
}
