package pong

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ ServerInterface = (*Server)(nil)

type Server struct{}

func NewServer() Server {
	return Server{}
}

// GetPing (GET /ping)
func (s Server) GetPing(ctx echo.Context) error {
	resp := Pong{
		Ping: "pong",
	}

	ctx.Logger().Infof("Ping response %s", resp.Ping)
	return ctx.JSON(http.StatusOK, resp)
}

// GetRestricted (GET /restricted)
func (s Server) GetRestricted(ctx echo.Context) error {
	name := s.printToken(ctx)

	resp := Pong{
		Ping: "Welcome " + name + "!",
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s Server) printToken(ctx echo.Context) string {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	ctx.Logger().Infof("Restricted name %s", name)
	return name
}
