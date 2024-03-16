package session

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/ucho456job/my_authn_authz/internal/constants"

	"github.com/labstack/echo/v4"
)

type SessionManager interface {
	SaveLoginSession(c echo.Context, userID string) error
}

type DefaultSessionManager struct{}

func NewDefaultSessionManager() SessionManager {
	return &DefaultSessionManager{}
}

func (dsm *DefaultSessionManager) SaveLoginSession(c echo.Context, userID string) error {
	session, err := session.Get(constants.LOGIN_SESSION_NAME, c)
	if err != nil {
		return err
	}

	session.Values["userId"] = userID
	session.Values["isLogin"] = true
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60,
		HttpOnly: true,
	}

	secure, err := strconv.ParseBool(os.Getenv("COOKIE_SECURE"))
	if err != nil {
		return err
	}

	if secure {
		session.Options.Secure = true
		session.Options.SameSite = http.SameSiteNoneMode
	}

	return session.Save(c.Request(), c.Response())
}
