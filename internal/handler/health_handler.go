package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/util/response"
)

type HealthHandler interface {
	Check(c echo.Context) error
}

type healthHandler struct {
	healthRepository repository.HealthRepository
	sessionManager   repository.SessionManagerRepository
}

func NewHealthHandler(
	healthRepository repository.HealthRepository,
	sessionManager repository.SessionManagerRepository,
) *healthHandler {
	return &healthHandler{
		healthRepository: healthRepository,
		sessionManager:   sessionManager,
	}
}

func (h *healthHandler) Check(c echo.Context) error {
	key := "healthCheckKey"

	postgresValue, err := h.healthRepository.CheckPostgres(c, key)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	redisValue, err := h.sessionManager.CheckRedis(c, key)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	config.Logger.Info("Server, Postgres and Redis are healthy", "postgresValue", postgresValue, "redisValue", redisValue)
	return c.String(http.StatusOK, "Server is healthy!")
}
