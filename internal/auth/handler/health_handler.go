package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/auth/repository"
	"github.com/ucho456job/lgtmeme/internal/util"
)

type HealthHandler interface {
	Check(c echo.Context) error
}

type healthHandler struct {
	healthRepository repository.HealthRepository
	sessionManager   repository.SessionManager
}

func NewHealthHandler(
	healthRepository repository.HealthRepository,
	sessionManager repository.SessionManager,
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
		return util.InternalServerErrorResponse(c, err)
	}

	redisValue, err := h.sessionManager.CheckRedis(c, key)
	if err != nil {
		return util.InternalServerErrorResponse(c, err)
	}

	config.Logger.Info("Server, Postgres and Redis are healthy", "postgresValue", postgresValue, "redisValue", redisValue)
	return c.String(http.StatusOK, "Server is healthy!")
}
