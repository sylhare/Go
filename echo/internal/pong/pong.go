package pong

import (
	"echo/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func PongEchoServer() *echo.Echo {
	server := NewServer()

	e := echo.New()
	e.Use(middleware.Logger)
	e.Use(middleware.JWT)
	e.Logger.SetLevel(log.DEBUG)
	RegisterHandlers(e, server)
	//r := e.Group("/restricted")
	//{
	//	r.Use(middleware.JWT)
	//	r.GET("", server.GetRestricted)
	//}
	return e
}
