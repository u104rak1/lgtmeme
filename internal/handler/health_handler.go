package handler

import (
	"net/http"
	"os"

	"golang.org/x/exp/slog"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/internal/repository"
	"github.com/ucho456job/my_authn_authz/internal/session"
)

type HealthHandler interface {
	CheckHealth(c echo.Context) error
}

type DefaultHealthHandler struct {
	userRepository repository.UserRepository
	sessionManager session.SessionManager
	logger         *slog.Logger
}

func NewHealthHandler(userRepository repository.UserRepository, sessionManager session.SessionManager, logger *slog.Logger) *DefaultHealthHandler {
	return &DefaultHealthHandler{
		userRepository: userRepository,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

func (h *DefaultHealthHandler) CheckHealth(c echo.Context) error {
	healthCheckKey := os.Getenv("HEALTH_CHECK_KEY")

	user, err := h.userRepository.FindByName(healthCheckKey)
	if err != nil {
		h.logger.Error("Failed to check helath for postgres", "error", err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
	}

	value, err := h.sessionManager.CheckHealthForRedis(healthCheckKey)
	if err != nil {
		h.logger.Error("Failed to check health for redis", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to check health for redis"})
	}

	h.logger.Info("Server, Postgres and Redis are healthy", "postgres_value", user.ID.String(), "redis_value", value)
	return c.String(http.StatusOK, "Server is healthy!")
}
