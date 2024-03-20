package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/config"
	"github.com/ucho456job/my_authn_authz/internal/handler"
	"github.com/ucho456job/my_authn_authz/internal/repository"
	"github.com/ucho456job/my_authn_authz/internal/util"
)

func main() {
	// Init config
	config.InitEnv()
	config.InitDB()
	config.InitSessionStore()
	config.InitLogger()

	// Init repository
	healthCheckRepo := repository.NewHealthCheckRepository(config.DB)
	oauthClientRepo := repository.NewOauthClientRepository(config.DB)
	userRepo := repository.NewUserRepository(config.DB)

	// Init session manager
	sessManager := util.NewSessionManager(config.Store, config.Pool)

	// Init handler
	authzHandler := handler.NewAuthorizationHandler(oauthClientRepo, userRepo, sessManager)
	healthHandler := handler.NewHealthHandler(healthCheckRepo, sessManager, config.Logger)
	loginHandler := handler.NewLoginHandler(userRepo, sessManager, config.Logger)

	e := echo.New()

	e.Use(config.SessionMiddleware(), config.LoggerMiddleware)

	e.Static("/", "view/out")
	e.GET(util.LOGIN_SCREEN_ENDPOINT, func(c echo.Context) error {
		return c.File(util.LOGIN_SCREEN_FILEPATH)
	})

	e.GET(util.AUTHORAIZETION_ENDPOINT, authzHandler.AuthorizationHandle)
	e.HEAD(util.HEALTH_ENDPOINT, healthHandler.CheckHealth)
	e.POST(util.LOGIN_ENDPOINT, loginHandler.Login)

	port := ":" + os.Getenv("PORT")
	e.Logger.Fatal(e.Start(port))
}
