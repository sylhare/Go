package main

import (
	"echo/internal/middleware"
	"echo/internal/pong"
	"log"
)

func main() {
	e := pong.PongEchoServer()
	e.Use(middleware.Logger)
	log.Fatal(e.Start("0.0.0.0:8080"))
}
