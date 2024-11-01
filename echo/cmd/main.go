package main

import (
	"echo/internal/pong"
	"log"
)

func main() {
	e := pong.PongEchoServer()
	log.Fatal(e.Start("0.0.0.0:8080"))
}
