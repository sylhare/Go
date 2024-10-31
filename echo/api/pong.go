package api

import (
	"github.com/labstack/echo/v4"
	"log"
)

func PongEchoServer() {
	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	server := NewServer()

	e := echo.New()

	RegisterHandlers(e, server)

	// And we serve HTTP until the world ends.
	log.Fatal(e.Start("0.0.0.0:8080"))

}
