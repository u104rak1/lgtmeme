package repository_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
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
		ClientID:    uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
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

func TestOauthClientRepository_FindByClientID(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	oc := model.OauthClient{
		ID:             uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Name:           "testName",
		ClientID:       uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
		ClientSecret:   "testSecret",
		RedirectURI:    "http://example.com/callback",
		ApplicationURL: "http://example.com",
		ClientType:     "confidential",
	}

	sqlStatement := `SELECT * FROM "oauth_clients" WHERE client_id = $1 ORDER BY "oauth_clients"."id" LIMIT $2`

	tests := []struct {
		name      string
		setupMock func()
		arg       func() (uuid.UUID, []string)
		want      *model.OauthClient
		isErr     bool
	}{
		{
			name: "Return oauth_client, do not specify columns",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "client_id", "client_secret", "redirect_uri", "application_url", "client_type"}).
					AddRow(oc.ID, oc.Name, oc.ClientID, oc.ClientSecret, oc.RedirectURI, oc.ApplicationURL, oc.ClientType)
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(oc.ClientID, 1).
					WillReturnRows(rows)
			},
			arg: func() (uuid.UUID, []string) {
				return oc.ClientID, []string{}
			},
			want:  &oc,
			isErr: false,
		},
		{
			name: "Return oauth_client, specify columns",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"client_id", "client_secret", "application_url"}).
					AddRow(oc.ClientID, oc.ClientSecret, oc.ApplicationURL)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "client_id","client_secret","application_url" FROM "oauth_clients" WHERE client_id = $1 ORDER BY "oauth_clients"."id" LIMIT $2`)).
					WithArgs(oc.ClientID, 1).
					WillReturnRows(rows)
			},
			arg: func() (uuid.UUID, []string) {
				return oc.ClientID, []string{"client_id", "client_secret", "application_url"}
			},
			want: &model.OauthClient{
				ClientID:       oc.ClientID,
				ClientSecret:   oc.ClientSecret,
				ApplicationURL: oc.ApplicationURL,
			},
			isErr: false,
		},
		{
			name: "Return error, because record not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(oc.ClientID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			arg: func() (uuid.UUID, []string) {
				return oc.ClientID, []string{}
			},
			want:  nil,
			isErr: true,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(oc.ClientID, 1).
					WillReturnError(testutil.ErrDB)
			},
			arg: func() (uuid.UUID, []string) {
				return oc.ClientID, []string{}
			},
			want:  nil,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewOauthClientRepository(db)
			clientID, columns := tt.arg()
			actual, err := repo.FirstByClientID(c, clientID, columns)

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
