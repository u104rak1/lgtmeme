package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	validator := config.InitValidator()

	// Init repository
	healthCheckRepo := repository.NewHealthCheckRepository(config.DB)
	oauthClientRepo := repository.NewOauthClientRepository(config.DB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(config.DB)
	sessManagerRepo := repository.NewSessionManager(config.Store, config.Pool)
	userRepo := repository.NewUserRepository(config.DB)

	// Init util
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	issuerURL := os.Getenv("BASE_URL")
	jwtService := util.NewJwtService(jwtKey, issuerURL)

	// Init handler
	authzHandler := handler.NewAuthorizationHandler(oauthClientRepo, userRepo, sessManagerRepo)
	healthHandler := handler.NewHealthHandler(healthCheckRepo, sessManagerRepo)
	loginHandler := handler.NewLoginHandler(userRepo, sessManagerRepo)
	tokenHandler := handler.NewTokenHandler(oauthClientRepo, refreshTokenRepo, userRepo, sessManagerRepo, jwtService)

	e := echo.New()

	e.Validator = validator

	e.Use(config.SessionMiddleware(), config.LoggerMiddleware)

	e.Static("/", "view/out")
	e.GET(util.LOGIN_SCREEN_ENDPOINT, func(c echo.Context) error {
		return c.File(util.LOGIN_SCREEN_FILEPATH)
	})
	e.GET(util.PASSKEY_SCREEN_ENDPOINT, func(c echo.Context) error {
		return c.File(util.PASSKEY_SCREEN_FILEPATH)
	})

	e.GET(util.AUTHORAIZETION_ENDPOINT, authzHandler.AuthorizationHandle)
	e.HEAD(util.HEALTH_ENDPOINT, healthHandler.CheckHealth)
	e.POST(util.LOGIN_ENDPOINT, loginHandler.Login)
	e.POST(util.TOKEN_ENDPOINT, tokenHandler.GenerateToken)

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
