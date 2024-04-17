package repository_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/util/timer"
	"github.com/ucho456job/lgtmeme/test/testutil"
)

func TestImageRepository(t *testing.T) {
	gormDB, mock := testutil.SetupMockDB(t)
	mockTimer := timer.MockTimer{}

	tUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	tURL := "http://example.com"
	tKeyword := "keyword"
	tUsedCount := 0
	tReported := false
	tConfirmed := false
	tCreatedAt := mockTimer.Now()

	sqlStatement := `INSERT INTO "images" ("url","keyword","used_count","reported","confirmed","id","created_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id","created_at"`

	createTests := []struct {
		name      string
		setupMock func()
		expectErr bool
	}{
		{
			name: "positive: Return nil, Create was successful",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(tURL, tKeyword, tUsedCount, tReported, tConfirmed, tUUID, tCreatedAt).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(tUUID, tCreatedAt))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "negative: Return error, because database connection error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(tURL, tKeyword, tUsedCount, tReported, tConfirmed, tUUID, tCreatedAt).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, tt := range createTests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(gormDB, &mockTimer)
			err := repo.Create(c, tUUID, tURL, tKeyword)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
