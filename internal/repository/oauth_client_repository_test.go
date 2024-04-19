package repository_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/lgtmeme/internal/dto"
	"github.com/ucho456job/lgtmeme/internal/model"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/test/testutil"
	"gorm.io/gorm"
)

func TestOauthClientRepository_IsValidOAuthClient(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	query := dto.AuthzQuery{
		ClientID:    testutil.TestUUIDs[0],
		RedirectURI: "http://example.com/callback",
		Scope:       "scope1 scope2",
	}

	scopes := []model.MasterScope{
		{Code: "scope1"},
		{Code: "scope2"},
	}

	sqlStatement := `SELECT osc.scope_code FROM oauth_clients AS oc INNER JOIN oauth_clients_scopes AS osc ON oc.id = osc.oauth_client_id WHERE oc.client_id = $1 AND oc.redirect_uri = $2`

	tests := []struct {
		name      string
		setupMock func()
		arg       func() dto.AuthzQuery
		want      bool
		isErr     bool
	}{
		{
			name: "Return true",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"scope_code"}).
					AddRow(scopes[0].Code).
					AddRow(scopes[1].Code)
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(query.ClientID, query.RedirectURI).
					WillReturnRows(rows)
			},
			arg: func() dto.AuthzQuery {
				return query
			},
			want:  true,
			isErr: false,
		},
		{
			name: "Return false, because the query has an invalid scope",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"scope_code"}).
					AddRow(scopes[0].Code).
					AddRow(scopes[1].Code)
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(query.ClientID, query.RedirectURI).
					WillReturnRows(rows)
			},
			arg: func() dto.AuthzQuery {
				return dto.AuthzQuery{
					ClientID:    query.ClientID,
					RedirectURI: query.RedirectURI,
					Scope:       "invalidScope",
				}
			},
			want:  false,
			isErr: false,
		},
		{
			name: "Return false, because record not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(query.ClientID, query.RedirectURI).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			arg: func() dto.AuthzQuery {
				return query
			},
			want:  false,
			isErr: false,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(query.ClientID, query.RedirectURI).
					WillReturnError(testutil.ErrDB)
			},
			arg: func() dto.AuthzQuery {
				return query
			},
			want:  false,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewOauthClientRepository(db)
			q := tt.arg()
			actual, err := repo.IsValidOAuthClient(c, q)

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

// TODO: Implement test
func TestOauthClientRepository_FindByClientIDWithScopes(t *testing.T) {
	t.Skip("Skip because sqlmock does not work well when using 'Preload'")
}
