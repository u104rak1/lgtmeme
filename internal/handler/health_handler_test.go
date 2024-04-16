package handler_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/lgtmeme/internal/handler"
	mock_repository "github.com/ucho456job/lgtmeme/test/mock/repository"
	"github.com/ucho456job/lgtmeme/test/testutil"
)

func TestHealthHandler_Check(t *testing.T) {
	testutil.SetupTestLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHealthRepo := mock_repository.NewMockHealthRepository(ctrl)
	mockSessionManager := mock_repository.NewMockSessionManagerRepository(ctrl)

	tests := []struct {
		name           string
		setupMocks     func(*gomock.Controller, echo.Context)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "positive: Return 200",
			setupMocks: func(ctrl *gomock.Controller, c echo.Context) {
				mockHealthRepo.EXPECT().CheckPostgres(c, "healthCheckKey").Return("OK", nil)
				mockSessionManager.EXPECT().CheckRedis(c, "healthCheckKey").Return("OK", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Server is healthy!",
		},
		{
			name: "negative: Return 500, because postgres error",
			setupMocks: func(ctrl *gomock.Controller, c echo.Context) {
				mockHealthRepo.EXPECT().CheckPostgres(c, "healthCheckKey").Return("", errors.New("postgres error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errorCode":"internal_server_error","errorMessage":"postgres error"}`,
		},
		{
			name: "negative: Return 500, because redis error",
			setupMocks: func(ctrl *gomock.Controller, c echo.Context) {
				mockHealthRepo.EXPECT().CheckPostgres(c, "healthCheckKey").Return("OK", nil)
				mockSessionManager.EXPECT().CheckRedis(c, "healthCheckKey").Return("", errors.New("redis error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errorCode":"internal_server_error","errorMessage":"redis error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := testutil.SetupMinEchoContext()
			handler := handler.NewHealthHandler(mockHealthRepo, mockSessionManager)
			tt.setupMocks(ctrl, c)
			err := handler.Check(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(rec.Body.String()))
		})
	}
}
