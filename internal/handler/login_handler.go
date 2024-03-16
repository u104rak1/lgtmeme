package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/internal/repository"
	"github.com/ucho456job/my_authn_authz/internal/session"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler interface {
	Login(c echo.Context)
}

type DefaultLoginHandler struct {
	userRepository repository.UserRepository
	sessionManager session.SessionManager
}

func NewSessionHandler(userRepository repository.UserRepository, sessionManager session.SessionManager) *DefaultLoginHandler {
	return &DefaultLoginHandler{
		userRepository: userRepository,
		sessionManager: sessionManager,
	}
}

func (h *DefaultLoginHandler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.userRepository.FindByName(username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
	}

	if err := h.sessionManager.SaveLoginSession(c, user.ID.String()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save session"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login success",
	})
}
