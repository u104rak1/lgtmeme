package repository_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/test/testutil"
	"gorm.io/gorm"
)

func TestCheckPostgres(t *testing.T) {
	gormDB, mock := testutil.SetupMockDB(t)

	sqlStatement := `SELECT "value" FROM "health_checks" WHERE key = $1 ORDER BY "health_checks"."key" LIMIT $2`

	tests := []struct {
		name          string
		key           string
		setupMock     func()
		expectedValue string
		expectErr     bool
	}{
		{
			name:          "positive: Return value",
			key:           "testKey",
			expectedValue: "testValue",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"value"}).AddRow("testValue")
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).WithArgs("testKey", 1).WillReturnRows(rows)
			},
			expectErr: false,
		},
		{
			name:          "negative: Return error, because record not found",
			key:           "missingKey",
			expectedValue: "",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).WithArgs("missingKey", 1).WillReturnError(gorm.ErrRecordNotFound)
			},
			expectErr: true,
		},
		{
			name:          "negative: Return error, because database connection error",
			key:           "anyKey",
			expectedValue: "",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).WithArgs("anyKey", 1).WillReturnError(errors.New("database connection failed"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewHealthRepository(gormDB)
			value, err := repo.CheckPostgres(c, tt.key)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, value)
			}
		})
	}
}
