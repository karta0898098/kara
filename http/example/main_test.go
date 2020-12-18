package main

import (
	"testing"

	"github.com/go-resty/resty/v2"
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
}

func (s *echoSuite) TearDownTest() {
}

func (s *echoSuite) TestRunEcho() {
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(`{"username":"testuser", "password":"testpass"}`)).
		Post("http://127.0.0.1:8080/ping")
	s.Equal(nil, err)
}
