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

func TestOauthClientRepository_ExistsForAuthz(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	tQuery := dto.AuthzQuery{
		ClientID:    uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		RedirectURI: "http://example.com/callback",
		Scope:       "scope1 scope2",
	}

	tScopes := []model.MasterScope{
		{Code: "scope1"},
		{Code: "scope2"},
	}

	sqlStatement := `SELECT osc.scope_code FROM oauth_clients oc INNER JOIN oauth_clients_scopes osc ON oc.id = osc.oauth_client_id WHERE oc.client_id = $1 AND oc.redirect_uri = $2`

	tests := []struct {
		name      string
		setupMock func()
		query     dto.AuthzQuery
		result    bool
		isErr     bool
	}{
		{
			name: "positive: Return true, client has all scopes",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"scope_code"}).
					AddRow(tScopes[0].Code).
					AddRow(tScopes[1].Code)
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(tQuery.ClientID, tQuery.RedirectURI).
					WillReturnRows(rows)
			},
			query:  tQuery,
			result: true,
			isErr:  false,
		},
		{
			name: "negative: Return false, because client doesn't have all scopes",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"scope_code"}).
					AddRow(tScopes[0].Code).
					AddRow(tScopes[1].Code)
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(tQuery.ClientID, tQuery.RedirectURI).
					WillReturnRows(rows)
			},
			query: dto.AuthzQuery{
				ClientID:    tQuery.ClientID,
				RedirectURI: tQuery.RedirectURI,
				Scope:       "invalidScope",
			},
			result: false,
			isErr:  false,
		},
		{
			name: "negative: Return false, because client doesn't have all scopes",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"scope_code"}).
					AddRow(tScopes[0].Code).
					AddRow(tScopes[1].Code)
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(tQuery.ClientID, tQuery.RedirectURI).
					WillReturnRows(rows)
			},
			query: dto.AuthzQuery{
				ClientID:    tQuery.ClientID,
				RedirectURI: tQuery.RedirectURI,
				Scope:       "invalidScope",
			},
			result: false,
			isErr:  false,
		},
		{
			name: "negative: Return false, because client not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(tQuery.ClientID, tQuery.RedirectURI).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			query:  tQuery,
			result: false,
			isErr:  false,
		},
		{
			name: "negative: Return error, because database connection error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(tQuery.ClientID, tQuery.RedirectURI).
					WillReturnError(testutil.ErrDBConnection)
			},
			query:  tQuery,
			result: false,
			isErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewOauthClientRepository(db)
			result, err := repo.ExistsForAuthz(c, tt.query)

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

func TestOauthClientRepository_FindByClientID(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	tClientID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	sqlStatement := `SELECT * FROM "oauth_clients" WHERE client_id = $1 ORDER BY "oauth_clients"."id" LIMIT $2`

	tests := []struct {
		name      string
		setupMock func()
		clientID  uuid.UUID
		result    *model.OauthClient
		isErr     bool
	}{
		{
			name: "positive: Return client",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "client_id", "client_secret", "redirect_uri", "application_url", "client_type"}).
					AddRow("123e4567-e89b-12d3-a456-426614174000", "testName", "123e4567-e89b-12d3-a456-426614174000", "testSecret", "http://example.com/callback", "http://example.com", "confidential")
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(tClientID, 1).
					WillReturnRows(rows)
			},
			clientID: tClientID,
			result: &model.OauthClient{
				ID:             tClientID,
				Name:           "testName",
				ClientID:       tClientID,
				ClientSecret:   "testSecret",
				RedirectURI:    "http://example.com/callback",
				ApplicationURL: "http://example.com",
				ClientType:     "confidential",
			},
			isErr: false,
		},
		{
			name: "negative: Return error, because client not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(tClientID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			clientID: tClientID,
			result:   nil,
			isErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewOauthClientRepository(db)
			result, err := repo.FindByClientID(c, tt.clientID)

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
