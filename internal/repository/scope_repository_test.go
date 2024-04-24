package repository_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/lgtmeme/internal/model"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/test/testutil"
)

func TestScopeRepositoryTest_FindByScopesStr(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	scopes := []model.MasterScope{
		{Code: "scope1", Description: "description1"},
		{Code: "scope2", Description: "description2"},
	}
	s1 := scopes[0]
	s2 := scopes[1]

	tests := []struct {
		name      string
		setupMock func()
		arg       func() string
		want      *[]model.MasterScope
		isErr     bool
	}{
		{
			name: "Return scopes",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"code", "description"}).
					AddRow(s1.Code, s1.Description).
					AddRow(s2.Code, s2.Description)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "master_scopes" WHERE code IN ($1,$2)`)).
					WithArgs(s1.Code, s2.Code).
					WillReturnRows(rows)
			},
			arg: func() string {
				return fmt.Sprintf("%s %s", s1.Code, s2.Code)
			},
			want:  &scopes,
			isErr: false,
		},
		{
			name: "Return empty",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"code", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "master_scopes" WHERE code IN ($1,$2)`)).
					WithArgs("invalidScope1", "invalidScope2").
					WillReturnRows(rows)
			},
			arg: func() string {
				return fmt.Sprintf("%s %s", "invalidScope1", "invalidScope2")
			},
			want:  &[]model.MasterScope{},
			isErr: false,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "master_scopes" WHERE code IN ($1,$2)`)).
					WithArgs(s1.Code, s2.Code).
					WillReturnError(testutil.ErrDB)
			},
			arg: func() string {
				return fmt.Sprintf("%s %s", s1.Code, s2.Code)
			},
			want:  nil,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			r := repository.NewScopeRepository(db)
			scopesStr := tt.arg()
			actual, err := r.FindByScopesStr(c, scopesStr)

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
