package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/auth/repository"
	"github.com/ucho456job/lgtmeme/internal/util"
)

type HealthHandler interface {
	CheckHealth(c echo.Context) error
}

type healthHandler struct {
	healthCheckRepository repository.HealthCheckRepository
	sessionManager        repository.SessionManager
}

func NewHealthHandler(
	healthCheckRepository repository.HealthCheckRepository,
	sessionManager repository.SessionManager,
) *healthHandler {
	return &healthHandler{
		healthCheckRepository: healthCheckRepository,
		sessionManager:        sessionManager,
	}
}

func (h *healthHandler) CheckHealth(c echo.Context) error {
	key := "healthCheckKey"

	postgresValue, err := h.healthCheckRepository.CheckHealthForPostgres(c, key)
	if err != nil {
		return util.InternalServerErrorResponse(c, err)
	}

	redisValue, err := h.sessionManager.CheckHealthForRedis(c, key)
	if err != nil {
		return util.InternalServerErrorResponse(c, err)
	}

	config.Logger.Info("Server, Postgres and Redis are healthy", "postgresValue", postgresValue, "redisValue", redisValue)
	return c.String(http.StatusOK, "Server is healthy!")
}
