package main

import (
	"os"

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
	config.InitLogger()

	// Init repository
	userRepo := repository.NewGormUserRepository(config.DB)

	// Init session manager
	sessionManager := session.NewDefaultSessionManager()

	// Init handler
	healthHandler := handler.NewHealthHandler(userRepo, sessionManager, config.Logger)
	loginHandler := handler.NewSessionHandler(userRepo, sessionManager)

	e := echo.New()

	e.Use(config.SessionMiddleware(), config.LoggerMiddleware)

	e.Static("/", "view/out")

	e.HEAD(constants.HEALTH_ENDPOINT, healthHandler.CheckHealth)
	e.POST(constants.LOGIN_ENDPOINT, loginHandler.Login)

	port := ":" + os.Getenv("PORT")
	e.Logger.Fatal(e.Start(port))
}
