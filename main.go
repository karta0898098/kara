package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/karta0898098/kara/db"
	"github.com/karta0898098/kara/http"
	"github.com/karta0898098/kara/zlog"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	Log      zlog.Config
	Database db.Config
	HTTP     http.Config
}

func (c *Config) New() Config {
	return *c
}

func SetRouter(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}

func main() {

	config := Config{
		Log: zlog.Config{
			Env:   "dev",
			AppID: "app",
			Debug: true,
			Local: true,
		},
		Database: db.Config{
			Read: db.Database{
				Debug:    false,
				Host:     "127.0.0.1",
				User:     "rode",
				Port:     3306,
				Password: "rode@3306",
				Name:     "dms",
				Type:     "mysql",
			},
			Write: db.Database{
				Debug:    false,
				Host:     "127.0.0.1",
				User:     "rode",
				Port:     3306,
				Password: "rode@3306",
				Name:     "dms",
				Type:     "mysql",
			},
		},
		HTTP: http.Config{
			Mode: "debug",
			Port: ":5500",
		},
	}
	app := fx.New(
		fx.Provide(
			config.New,
			config.Log.New,
			config.Database.New,
			config.HTTP.New,
			db.NewConnection,
			http.RunGin,
		),
		fx.Invoke(zlog.Setup),
		fx.Invoke(SetRouter),
	)

	exitCode := 0
	if err := app.Start(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(exitCode)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stop

	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		fmt.Println(err)
	}

	os.Exit(exitCode)
}
