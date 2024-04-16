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
	e := testutil.BeforeAll(t)
	defer testutil.AfterAll(t)

	tests := []struct {
		name           string
		prepare        func(t *testing.T)
		expectedStatus int
		expectedBody   string
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
			expectedBody:   "Server is healthy!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(t)
			req := httptest.NewRequest(http.MethodHead, config.HEALTH_ENDPOINT, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.expectedBody, rec.Body.String())
		})
	}
}
