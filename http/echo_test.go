package http

import (
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/bwmarrin/snowflake"
	"github.com/karta0898098/kara/zlog"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type echoSuite struct {
	suite.Suite
	app    *fx.App
	engine *echo.Echo
}

func TestEndpoint(t *testing.T) {
	suite.Run(t, new(echoSuite))
}

func (s *echoSuite) SetupTest() {
	node, _ := snowflake.NewNode(1)
	httpConfig := &Config{
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

	zlog.Setup(logConfig)
	s.engine = NewEcho(httpConfig, node)
	s.engine.GET("/hello", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"data": "hello world",
		})
	})
}

func (s *echoSuite) TearDownTest() {
}

func (s *echoSuite) TestNewEcho() {
	r := gofight.New()
	r.GET("/hello").
		SetDebug(false).
		Run(s.engine, func(resp gofight.HTTPResponse, req gofight.HTTPRequest) {
			s.Equal(200, resp.Code)
		})
}
