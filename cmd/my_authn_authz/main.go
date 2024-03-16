package main

import (
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/config"
	"github.com/ucho456job/my_authn_authz/internal/constants"
	"github.com/ucho456job/my_authn_authz/internal/handler"
	"github.com/ucho456job/my_authn_authz/internal/repository"
	"github.com/ucho456job/my_authn_authz/internal/session"
)

func main() {
	// Init config
	config.InitEnv()
	config.InitDB()
	config.InitSessionStore()

	// Init repository
	userRepo := repository.NewGormUserRepository(config.DB)

	// Init session manager
	sessionManager := session.NewDefaultSessionManager()

	// Init handler
	loginHandler := handler.NewSessionHandler(userRepo, sessionManager)

	e := echo.New()

	e.Use(config.SessionMiddleware())

	e.Static("/", "view/out")

	e.POST(constants.LOGIN_ENDPOINT, loginHandler.Login)

	e.Logger.Fatal(e.Start(":8080"))
}
