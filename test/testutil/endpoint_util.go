package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	"github.com/sebdah/goldie/v2"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/model"
	"github.com/ucho456job/lgtmeme/internal/setup"
)

type RedisData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type HTTPReq struct {
	URL    string      `json:"url"`
	Method string      `json:"method"`
	Header interface{} `json:"header"`
}

type HTTPResp struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

type TestResult struct {
	BeforeDB    map[string]interface{} `json:"beforeDB"`
	AfterDB     map[string]interface{} `json:"afterDB"`
	BeforeRedis []RedisData            `json:"beforeRedis"`
	AfterRedis  []RedisData            `json:"afterRedis"`
	Request     HTTPReq                `json:"request"`
	Response    HTTPResp               `json:"response"`
}

func BeforeAll(t *testing.T, folderName string) (*echo.Echo, *goldie.Goldie) {
	// Change current directory to root directory
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	os.Chdir(filepath.Join(basepath, "../.."))

	t.Setenv("ECHO_MODE", "local")
	t.Setenv("LOG_LEVEL", "SILENT")
	e := setup.SetupServer()

	ClearAllData(t)

	goldieDir := filepath.Join("test", "endpoint", "testdata", folderName)
	gol := goldie.New(t, goldie.WithFixtureDir(goldieDir))

	return e, gol
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
		if err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
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

func ClearAllData(t *testing.T) {
	ClearDB(t)
	ClearRedisData(t)
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

func FetchDBData(t *testing.T, tableNames []string) map[string]interface{} {
	db := config.DB
	allData := make(map[string]interface{})

	for _, tableName := range tableNames {
		modelType := getModelType(tableName)
		if modelType == nil {
			continue
		}
		modelSlice := reflect.New(reflect.SliceOf(modelType)).Interface()

		if err := db.Table(tableName).Find(modelSlice).Error; err != nil {
			t.Fatal(err)
		}

		allData[tableName] = modelSlice
	}

	return allData
}

func getModelType(tableName string) reflect.Type {
	switch tableName {
	case "health_checks":
		return reflect.TypeOf(&model.HealthCheck{})
	case "users":
		return reflect.TypeOf(&model.User{})
	case "oauth_clients":
		return reflect.TypeOf(&model.OauthClient{})
	case "master_scopes":
		return reflect.TypeOf(&model.MasterScope{})
	case "oauth_clients_scopes":
		return reflect.TypeOf(&model.OauthClientsScopes{})
	case "master_application_types":
		return reflect.TypeOf(&model.MasterApplicationType{})
	case "oauth_clients_application_types":
		return reflect.TypeOf(&model.OauthClientsApplicationTypes{})
	case "refresh_tokens":
		return reflect.TypeOf(&model.RefreshToken{})
	case "images":
		return reflect.TypeOf(&model.Image{})
	default:
		return nil
	}
}

func FetchRedisData(t *testing.T) []RedisData {
	conn := config.Pool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		t.Fatal(err)
	}

	var data []RedisData
	for _, key := range keys {
		value, err := redis.String(conn.Do("GET", key))
		if err != nil {
			t.Fatal(err)
		}
		data = append(data, RedisData{Key: key, Value: value})
	}

	return data
}

func GenerateResultJSON(
	t *testing.T,
	beforeDBData map[string]interface{},
	afterDBData map[string]interface{},
	beforeRedisData []RedisData,
	afterRedisData []RedisData,
	req *http.Request,
	rec *httptest.ResponseRecorder,
) []byte {
	result := TestResult{
		BeforeDB:    beforeDBData,
		AfterDB:     afterDBData,
		BeforeRedis: beforeRedisData,
		AfterRedis:  afterRedisData,
		Request: HTTPReq{
			URL:    req.URL.String(),
			Method: req.Method,
			Header: req.Header,
		},
		Response: HTTPResp{
			StatusCode: rec.Code,
			Body:       rec.Body,
		},
	}

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal result to JSON: %v", err)
	}
	return resultJSON
}
