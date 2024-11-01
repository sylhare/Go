package pong

import (
	"github.com/labstack/echo/v4"
)

func PongEchoServer() *echo.Echo {
	server := NewServer()

	e := echo.New()

	RegisterHandlers(e, server)
	return e
}
