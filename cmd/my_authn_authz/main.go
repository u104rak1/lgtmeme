package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/config"
	"github.com/ucho456job/my_authn_authz/internal/constant"
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
	userRepo := repository.NewUserRepository(config.DB)

	// Init session manager
	sessionManager := session.NewSessionManager()

	// Init handler
	healthHandler := handler.NewHealthHandler(userRepo, sessionManager, config.Logger)
	loginHandler := handler.NewLoginHandler(userRepo, sessionManager, config.Logger)

	e := echo.New()

	e.Use(config.SessionMiddleware(), config.LoggerMiddleware)

	e.Static("/", "view/out")

	e.HEAD(constant.HEALTH_ENDPOINT, healthHandler.CheckHealth)
	e.POST(constant.LOGIN_ENDPOINT, loginHandler.Login)

	port := ":" + os.Getenv("PORT")
	e.Logger.Fatal(e.Start(port))
}
