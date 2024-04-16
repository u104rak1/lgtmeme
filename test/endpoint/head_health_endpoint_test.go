package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/model"
	"github.com/ucho456job/lgtmeme/test/testutil"
)

func TestHealthCheckHandler(t *testing.T) {
	e, gol := testutil.BeforeAll(t, "HealthCheckHead")
	defer testutil.AfterAll(t)

	tests := []struct {
		name           string
		prepare        func(t *testing.T)
		expectedStatus int
	}{
		{
			name: "positive: Return 200",
			prepare: func(t *testing.T) {
				key := "healthCheckKey"
				prepareHealthCheck := model.HealthCheck{
					Key:   key,
					Value: "healthCheckValue",
				}
				testutil.PrepareDBData(t, &prepareHealthCheck)
				testutil.PrepareRedisData(t, key, "healthCheckValue")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "negative: Return 500, because health_checks not found",
			prepare:        func(t *testing.T) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "negative: Return 500, because redis value not found",
			prepare: func(t *testing.T) {
				key := "healthCheckKey"
				prepareHealthCheck := model.HealthCheck{
					Key:   key,
					Value: "healthCheckValue",
				}
				testutil.PrepareDBData(t, &prepareHealthCheck)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(t)
			beforeDBData := testutil.FetchDBData(t, []string{"health_checks"})
			beforeRedisData := testutil.FetchRedisData(t)

			req := httptest.NewRequest(http.MethodHead, config.HEALTH_ENDPOINT, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			afterDBData := testutil.FetchDBData(t, []string{"health_checks"})
			afterRedisData := testutil.FetchRedisData(t)

			resultJSON := testutil.GenerateResultJSON(t, beforeDBData, afterDBData, beforeRedisData, afterRedisData, req, rec)
			gol.Assert(t, tt.name, resultJSON)
			testutil.ClearAllData(t)
		})
	}
}
