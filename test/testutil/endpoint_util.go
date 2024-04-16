package testutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/model"
	"github.com/ucho456job/lgtmeme/internal/setup"
)

func BeforeAll(t *testing.T) *echo.Echo {
	// Change current directory to root directory
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	os.Chdir(filepath.Join(basepath, "../.."))

	t.Setenv("ECHO_MODE", "local")
	t.Setenv("LOG_LEVEL", "SILENT")
	e := setup.SetupServer()

	ClearDB(t)
	ClearRedisData(t)

	return e
}

func AfterAll(t *testing.T) {
	setup.CloseConnection()
}

func ClearDB(t *testing.T) {
	db := config.DB

	tables := []string{
		"images",
		"health_checks",
		"refresh_tokens",
		"oauth_clients_application_types",
		"oauth_clients_scopes",
		"oauth_clients",
		"users",
		"master_application_types",
		"master_scopes",
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			t.Fatal(err)
		}
		if err := db.AutoMigrate(
			&model.HealthCheck{},
			&model.User{},
			&model.OauthClient{},
			&model.MasterScope{},
			&model.OauthClientsScopes{},
			&model.MasterApplicationType{},
			&model.OauthClientsApplicationTypes{},
			&model.RefreshToken{},
			&model.Image{},
		); err != nil {
			t.Fatal(err)
		}
	}
}

func ClearRedisData(t *testing.T) {
	conn := config.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("FLUSHDB")
	if err != nil {
		t.Fatal(err)
	}
}

func PrepareDBData[T any](t *testing.T, data T) {
	db := config.DB
	if err := db.Create(&data).Error; err != nil {
		t.Fatal(err)
	}
}

func PrepareRedisData(t *testing.T, key, value string) {
	conn := config.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		t.Fatal(err)
	}
}
