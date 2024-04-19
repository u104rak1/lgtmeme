package repository_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/test/testutil"
	"gorm.io/gorm"
)

func TestHealthRepository_CheckPostgres(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	sqlStatement := `SELECT "value" FROM "health_checks" WHERE key = $1 ORDER BY "health_checks"."key" LIMIT $2`
	ckey := "testKey"

	tests := []struct {
		name      string
		setupMock func()
		arg       func() string
		want      string
		isErr     bool
	}{
		{
			name: "Return value",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"value"}).AddRow("testValue")
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(ckey, 1).
					WillReturnRows(rows)
			},
			arg: func() string {
				return ckey
			},
			want:  "testValue",
			isErr: false,
		},
		{
			name: "Return error, because record not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs("missingKey", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			arg: func() string {
				return "missingKey"
			},
			want:  "",
			isErr: true,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(ckey, 1).
					WillReturnError(testutil.ErrDB)
			},
			arg: func() string {
				return ckey
			},
			want:  "",
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewHealthRepository(db)
			key := tt.arg()
			actual, err := repo.CheckPostgres(c, key)

			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, actual)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
