package testutil

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	return gormDB, mock
}

func SetupMinEchoContext() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetResponse(echo.NewResponse(rec, e))
	return c, rec
}

func SetupTestLogger() {
	config.Logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
}
