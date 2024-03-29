package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	authHander "github.com/ucho456job/lgtmeme/internal/auth/handler"
	authRepository "github.com/ucho456job/lgtmeme/internal/auth/repository"
	authService "github.com/ucho456job/lgtmeme/internal/auth/service"
	clientHandler "github.com/ucho456job/lgtmeme/internal/client/handler"
	clientRepository "github.com/ucho456job/lgtmeme/internal/client/repository"
	clientService "github.com/ucho456job/lgtmeme/internal/client/service"
)

func main() {
	// Init Config
	config.InitEnv()
	config.InitDB()
	config.InitSessionStore()
	config.InitLogger()
	validator := config.InitValidator()

	// Init Auth Repository
	healthCheckRepository := authRepository.NewHealthCheckRepository(config.DB)
	oauthClientRepository := authRepository.NewOauthClientRepository(config.DB)
	refreshTokenRepository := authRepository.NewRefreshTokenRepository(config.DB)
	authSessionManagerRepository := authRepository.NewSessionManager(config.Store, config.Pool)
	userRepository := authRepository.NewUserRepository(config.DB)

	// Init Client Repository
	clientSessionManagerRepository := clientRepository.NewSessionManager(config.Store, config.Pool)

	// Init Client Service
	clientCredentialsService := clientService.NewClientCredentialsService()

	// Init Auth Service
	jwtService := authService.NewJwtService()

	// Init Echo
	e := echo.New()
	e.Validator = validator
	e.Use(config.SessionMiddleware(), config.LoggerMiddleware)
	e.Static(config.STATIC_ENDPOINT, config.STATIC_FILEPATH)

	// Init Auth Handler
	authzHandler := authHander.NewAuthorizationHandler(oauthClientRepository, userRepository, authSessionManagerRepository)
	healthHandler := authHander.NewHealthHandler(healthCheckRepository, authSessionManagerRepository)
	jwksHandler := authHander.NewJwksHandler(jwtService)
	loginHandler := authHander.NewLoginHandler(userRepository, authSessionManagerRepository)
	logoutHandler := authHander.NewLogoutHandler(authSessionManagerRepository)
	tokenHandler := authHander.NewTokenHandler(oauthClientRepository, refreshTokenRepository, userRepository, authSessionManagerRepository, jwtService)
	e.GET(config.AUTHORAIZETION_ENDPOINT, authzHandler.AuthorizationHandle)
	e.HEAD(config.HEALTH_ENDPOINT, healthHandler.CheckHealth)
	e.GET(config.JWKS_ENDPOINT, jwksHandler.GetJwks)
	e.POST(config.LOGIN_ENDPOINT, loginHandler.Login)
	e.GET(config.LOGIN_VIEW_ENDPOINT, loginHandler.GetLoginView)
	e.GET(config.LOGOUT_ENDPOINT, logoutHandler.Logout)
	e.POST(config.TOKEN_ENDPOINT, tokenHandler.GenerateToken)

	// Init Resource Handler

	// Init Client Handler
	clientAuthHandler := clientHandler.NewAuthorizationHandler()
	errorViewHandler := clientHandler.NewErrorViewHandler()
	homeViewHandler := clientHandler.NewHomeViewHandler(clientSessionManagerRepository, clientCredentialsService)
	e.GET(config.CLIENT_AUTH_ENDPOINT, clientAuthHandler.RedirectAuthz)
	e.GET(config.ERROR_VIEW_ENDPOINT, errorViewHandler.GetErrorView)
	e.GET(config.HOME_VIEW_ENDPOINT, homeViewHandler.GetHomeView)

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
