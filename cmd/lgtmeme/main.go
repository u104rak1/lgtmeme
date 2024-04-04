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
	resourceHandler "github.com/ucho456job/lgtmeme/internal/resource/handler"
	"github.com/ucho456job/lgtmeme/internal/resource/middleware"
	resourceRepository "github.com/ucho456job/lgtmeme/internal/resource/repository"
	resourceService "github.com/ucho456job/lgtmeme/internal/resource/service"
	"github.com/ucho456job/lgtmeme/internal/util/clock"
	"github.com/ucho456job/lgtmeme/internal/util/uuidgen"
)

func main() {
	config.NewEnv()
	config.NewDB()
	config.NewSessionStore()
	config.NewLogger()
	validator := config.NewValidator()

	e := echo.New()
	e.Validator = validator
	e.Use(config.SessionMiddleware(), config.LoggerMiddleware)
	e.Static(config.STATIC_ENDPOINT, config.STATIC_FILEPATH)

	newAuthServer(e)
	newClientServer(e)
	newResourceServer(e)

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
	healthRepo := authRepository.NewHealthRepository(config.DB)
	oauthClientRepo := authRepository.NewOauthClientRepository(config.DB)
	refreshTokenRepo := authRepository.NewRefreshTokenRepository(config.DB)
	authSessManaRepo := authRepository.NewSessionManager(config.Store, config.Pool)
	userRepo := authRepository.NewUserRepository(config.DB)

	jwtServ := authService.NewJWTService()

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
	sessManaRepo := clientRepository.NewSessionManager(config.Store, config.Pool)

	generalAccessTokenServ := clientService.NewGeneralAccessTokenService()
	imgServ := clientService.NewImageService()
	ownerAccessTokenServ := clientService.NewOwnerAccessTokenService()

	authHandler := clientHandler.NewAuthzHandler(sessManaRepo, ownerAccessTokenServ)
	errHandler := clientHandler.NewErrHandler()
	homeHandler := clientHandler.NewHomeHandler(sessManaRepo, generalAccessTokenServ)
	imgHandler := clientHandler.NewImageHandler(sessManaRepo, imgServ)

	e.GET(config.AUTH_VIEW_ENDPOINT, authHandler.GetView)
	e.GET(config.CLIENT_AUTH_ENDPOINT, authHandler.RedirectAuthz)
	e.GET(config.CLIENT_AUTH_CALLBACK_ENDPOINT, authHandler.Callback)

	e.GET(config.ERROR_VIEW_ENDPOINT, errHandler.GetView)

	e.GET(config.HOME_VIEW_ENDPOINT, homeHandler.GetView)

	e.GET(config.CREATE_IMAGE_VIEW_ENDPOINT, imgHandler.GetCreateImageView)
	e.POST(config.CLIENT_IMAGES_ENDPOINT, imgHandler.Post)
	e.GET(config.CLIENT_IMAGES_ENDPOINT, imgHandler.BulkGet)
	e.PATCH(config.CLIENT_IMAGES_ENDPOINT+"/:image_id", imgHandler.Patch)
}

func newResourceServer(e *echo.Echo) {
	imgRepo := resourceRepository.NewImageRepository(config.DB, &clock.RealClocker{})

	storageServ := resourceService.NewStorageService()

	imgHandler := resourceHandler.NewImageHandler(imgRepo, storageServ, &uuidgen.RealUUIDGenerator{})

	e.POST(config.RESOURCE_IMAGES_ENDPOINT, imgHandler.Post, middleware.VerifyAccessToken(config.IMAGES_CREATE_SCOPE))
	e.GET(config.RESOURCE_IMAGES_ENDPOINT, imgHandler.BulkGet, middleware.VerifyAccessToken(config.IMAGES_READ_SCOPE))
	e.PATCH(config.RESOURCE_IMAGES_ENDPOINT+"/:image_id", imgHandler.Patch, middleware.VerifyAccessToken(config.IMAGES_UPDATE_SCOPE))
}
