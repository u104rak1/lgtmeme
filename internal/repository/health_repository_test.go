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
		name      string
		setupMock func()
		key       string
		result    string
		isErr     bool
	}{
		{
			name: "positive: Return value",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"value"}).AddRow("testValue")
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).WithArgs("testKey", 1).WillReturnRows(rows)
			},
			key:    "testKey",
			result: "testValue",
			isErr:  false,
		},
		{
			name: "negative: Return error, because record not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).WithArgs("missingKey", 1).WillReturnError(gorm.ErrRecordNotFound)
			},
			key:    "missingKey",
			result: "",
			isErr:  true,
		},
		{
			name: "negative: Return error, because database connection error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).WithArgs("anyKey", 1).WillReturnError(errors.New("database connection failed"))
			},
			key:    "anyKey",
			result: "",
			isErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewHealthRepository(gormDB)
			result, err := repo.CheckPostgres(c, tt.key)

			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.result, result)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
