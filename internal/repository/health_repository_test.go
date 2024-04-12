package repository_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/testutil"
	"gorm.io/gorm"
)

func TestCheckPostgres(t *testing.T) {
	gormDB, mock := testutil.SetupMockDB(t)

	expectQuery := `SELECT "value" FROM "health_checks" WHERE key = \$1 ORDER BY "health_checks"\."key" LIMIT \$2`

	tests := []struct {
		name          string
		key           string
		setupMock     func()
		expectedValue string
		expectError   bool
	}{
		{
			name:          "positive: Return value",
			key:           "testKey",
			expectedValue: "testValue",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"value"}).AddRow("testValue")
				mock.ExpectQuery(expectQuery).WithArgs("testKey", 1).WillReturnRows(rows)
			},
			expectError: false,
		},
		{
			name:          "negative: Return error, because record not found",
			key:           "missingKey",
			expectedValue: "",
			setupMock: func() {
				mock.ExpectQuery(expectQuery).WithArgs("missingKey", 1).WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name:          "negative: Return error, because database connection error",
			key:           "anyKey",
			expectedValue: "",
			setupMock: func() {
				mock.ExpectQuery(expectQuery).WithArgs("anyKey", 1).WillReturnError(errors.New("database connection failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewHealthRepository(gormDB)
			value, err := repo.CheckPostgres(c, tt.key)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, value)
			}
		})
	}
}
