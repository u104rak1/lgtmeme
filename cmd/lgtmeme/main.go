package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/handler"
	authHandler "github.com/ucho456job/lgtmeme/internal/handler/auth"
	viewHandler "github.com/ucho456job/lgtmeme/internal/handler/view"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/util"
)

func main() {
	// Init config
	config.InitEnv()
	config.InitDB()
	config.InitSessionStore()
	config.InitLogger()
	validator := config.InitValidator()

	// Init repository
	healthCheckRepo := repository.NewHealthCheckRepository(config.DB)
	oauthClientRepo := repository.NewOauthClientRepository(config.DB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(config.DB)
	sessManagerRepo := repository.NewSessionManager(config.Store, config.Pool)
	userRepo := repository.NewUserRepository(config.DB)

	// Init util
	jwtService := util.NewJwtService()

	// Init echo
	e := echo.New()
	e.Validator = validator
	e.Use(config.SessionMiddleware(), config.LoggerMiddleware)

	// Init Health handler
	healthHandler := handler.NewHealthHandler(healthCheckRepo, sessManagerRepo)
	e.HEAD(config.HEALTH_ENDPOINT, healthHandler.CheckHealth)

	// Init Auth handler
	authzHandler := authHandler.NewAuthorizationHandler(oauthClientRepo, userRepo, sessManagerRepo)
	jwksHandler := authHandler.NewJwksHandler(jwtService)
	loginHandler := authHandler.NewLoginHandler(userRepo, sessManagerRepo)
	logoutHandler := authHandler.NewLogoutHandler(sessManagerRepo)
	tokenHandler := authHandler.NewTokenHandler(oauthClientRepo, refreshTokenRepo, userRepo, sessManagerRepo, jwtService)
	e.GET(config.AUTHORAIZETION_ENDPOINT, authzHandler.AuthorizationHandle)
	e.GET(config.JWKS_ENDPOINT, jwksHandler.GetJwks)
	e.POST(config.LOGIN_ENDPOINT, loginHandler.Login)
	e.GET(config.LOGOUT_ENDPOINT, logoutHandler.Logout)
	e.POST(config.TOKEN_ENDPOINT, tokenHandler.GenerateToken)

	// Init Resource handler

	// Init View handler
	loginViewHandler := viewHandler.NewLoginViewHandler()
	e.GET(config.LOGIN_VIEW_ENDPOINT, loginViewHandler.GetLoginView)
	e.GET(config.PASSKEY_VIEW_ENDPOINT, func(c echo.Context) error {
		return c.File(config.PASSKEY_VIEW_FILEPATH)
	})
	e.Static(config.STATIC_ENDPOINT, config.STATIC_FILEPATH)

	// Init API handler

	// Graceful shutdown
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	port := ":" + os.Getenv("PORT")
	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
