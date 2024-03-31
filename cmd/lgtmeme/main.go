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
	// Init config
	config.NewEnv()
	config.NewDB()
	config.NewSessionStore()
	config.NewLogger()
	validator := config.NewValidator()

	// Init echo
	e := echo.New()
	e.Validator = validator
	e.Use(config.SessionMiddleware(), config.LoggerMiddleware)
	e.Static(config.STATIC_ENDPOINT, config.STATIC_FILEPATH)

	newAuthServer(e)
	newClientServer(e)

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

func newAuthServer(e *echo.Echo) {
	// Init repository
	healthRepo := authRepository.NewHealthRepository(config.DB)
	oauthClientRepo := authRepository.NewOauthClientRepository(config.DB)
	refreshTokenRepo := authRepository.NewRefreshTokenRepository(config.DB)
	authSessManaRepo := authRepository.NewSessionManager(config.Store, config.Pool)
	userRepo := authRepository.NewUserRepository(config.DB)

	// Init service
	jwtServ := authService.NewJwtService()

	// Init handler
	authzHandler := authHander.NewAuthzHandler(oauthClientRepo, userRepo, authSessManaRepo)
	healthHandler := authHander.NewHealthHandler(healthRepo, authSessManaRepo)
	jwksHandler := authHander.NewJwksHandler(jwtServ)
	loginHandler := authHander.NewLoginHandler(userRepo, authSessManaRepo)
	logoutHandler := authHander.NewLogoutHandler(authSessManaRepo)
	tokenHandler := authHander.NewTokenHandler(oauthClientRepo, refreshTokenRepo, userRepo, authSessManaRepo, jwtServ)
	e.GET(config.AUTHZ_ENDPOINT, authzHandler.Authorize)
	e.HEAD(config.HEALTH_ENDPOINT, healthHandler.Check)
	e.GET(config.JWKS_ENDPOINT, jwksHandler.Get)
	e.POST(config.LOGIN_ENDPOINT, loginHandler.Login)
	e.GET(config.LOGIN_VIEW_ENDPOINT, loginHandler.GetView)
	e.GET(config.LOGOUT_ENDPOINT, logoutHandler.Logout)
	e.POST(config.TOKEN_ENDPOINT, tokenHandler.Generate)
}

func newClientServer(e *echo.Echo) {
	// Init repository
	clientSessManaRepo := clientRepository.NewSessionManager(config.Store, config.Pool)

	// Init service
	generalAccessTokenServ := clientService.NewGeneralAccessTokenService()
	ownerAccessTokenServ := clientService.NewOwnerAccessTokenService()

	// Init Handler
	clientAuthHandler := clientHandler.NewAuthzHandler(ownerAccessTokenServ)
	errHandler := clientHandler.NewErrHandler()
	homeHandler := clientHandler.NewHomeHandler(clientSessManaRepo, generalAccessTokenServ)
	e.GET(config.CLIENT_AUTH_ENDPOINT, clientAuthHandler.RedirectAuthz)
	e.GET(config.ERROR_VIEW_ENDPOINT, errHandler.GetView)
	e.GET(config.HOME_VIEW_ENDPOINT, homeHandler.GetView)
}
