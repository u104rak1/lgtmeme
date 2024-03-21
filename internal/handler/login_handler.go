package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/internal/repository"
	"github.com/ucho456job/my_authn_authz/internal/util"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type LoginHandler interface {
	Login(c echo.Context)
}

type loginHandler struct {
	userRepository repository.UserRepository
	sessionManager util.SessionManager
	logger         *slog.Logger
}

func NewLoginHandler(userRepository repository.UserRepository, sessionManager util.SessionManager, logger *slog.Logger) *loginHandler {
	return &loginHandler{
		userRepository: userRepository,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

func (h *loginHandler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.userRepository.FindByName(c, username)
	if err != nil {
		h.logger.Warn("Failed to find user", "error", err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		h.logger.Warn("Failed to compare password", "error", err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
	}

	if err := h.sessionManager.CacheLoginSession(c, user.ID.String()); err != nil {
		h.logger.Error("Failed to save login session", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save session"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login success",
	})
}
