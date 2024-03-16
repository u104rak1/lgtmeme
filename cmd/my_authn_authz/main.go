package main

import (
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/config"
)

func main() {
	e := echo.New()

	e.Use(config.SessionMiddleware())

	e.Static("/", "view/out")

	e.Logger.Fatal(e.Start(":8080"))
}
