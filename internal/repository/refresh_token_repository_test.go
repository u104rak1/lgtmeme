package repository_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/lgtmeme/internal/model"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/test/testutil"
	"gorm.io/gorm"
)

var (
	rt = model.RefreshToken{
		Token:    "token",
		UserID:   testutil.TestUUIDs[0],
		ClientID: testutil.TestUUIDs[1],
		Scopes:   "scope1 scope2",
	}
)

func TestRefreshTokenRepository_Create(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	sqlStatement := `INSERT INTO "refresh_tokens" ("token","user_id","client_id","scopes") VALUES ($1,$2,$3,$4)`

	tests := []struct {
		name      string
		setupMock func()
		arg       func() (uuid.UUID, uuid.UUID, string, string)
		isErr     bool
	}{
		{
			name: "Create was successful",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(sqlStatement)).
					WithArgs(rt.Token, rt.UserID, rt.ClientID, rt.Scopes).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			arg: func() (uuid.UUID, uuid.UUID, string, string) {
				return rt.UserID, rt.ClientID, rt.Token, rt.Scopes
			},
			isErr: false,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(sqlStatement)).
					WithArgs(rt.Token, rt.UserID, rt.ClientID, rt.Scopes).
					WillReturnError(testutil.ErrDB)
				mock.ExpectRollback()
			},
			arg: func() (uuid.UUID, uuid.UUID, string, string) {
				return rt.UserID, rt.ClientID, rt.Token, rt.Scopes
			},
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewRefreshTokenRepository(db)
			userID, clientID, token, scope := tt.arg()
			err := repo.Create(c, userID, clientID, token, scope)

			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRefreshTokenRepository_FirstByToken(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	sqlStatement := `SELECT * FROM "refresh_tokens" WHERE token = $1 ORDER BY "refresh_tokens"."token" LIMIT $2`

	tests := []struct {
		name      string
		setupMock func()
		arg       func() string
		want      *model.RefreshToken
		isErr     bool
	}{
		{
			name: "Return RefreshToken",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"token", "user_id", "client_id", "scopes"}).
					AddRow(rt.Token, rt.UserID, rt.ClientID, rt.Scopes)
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(rt.Token, 1).
					WillReturnRows(rows)
			},
			arg: func() string {
				return rt.Token
			},
			want:  &rt,
			isErr: false,
		},
		{
			name: "Return error, because record not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(rt.Token, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			arg: func() string {
				return rt.Token
			},
			want:  nil,
			isErr: true,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(rt.Token, 1).
					WillReturnError(testutil.ErrDB)
			},
			arg: func() string {
				return rt.Token
			},
			want:  nil,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewRefreshTokenRepository(db)
			token := tt.arg()
			actual, err := repo.FirstByToken(c, token)

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

func TestRerfeshTokenRepository_Update(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	tests := []struct {
		name      string
		setupMock func()
		arg       func() (uuid.UUID, uuid.UUID, string, string)
		isErr     bool
	}{
		{
			name: "Update was successful",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "refresh_tokens" WHERE user_id = $1 AND client_id = $2`)).
					WithArgs(rt.UserID, rt.ClientID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "refresh_tokens" ("token","user_id","client_id","scopes") VALUES ($1,$2,$3,$4)`)).
					WithArgs(rt.Token, rt.UserID, rt.ClientID, rt.Scopes).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			arg: func() (uuid.UUID, uuid.UUID, string, string) {
				return rt.UserID, rt.ClientID, rt.Token, rt.Scopes
			},
			isErr: false,
		},
		{
			name: "Return error, because record not found",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "refresh_tokens" WHERE user_id = $1 AND client_id = $2`)).
					WithArgs(rt.UserID, rt.ClientID).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			arg: func() (uuid.UUID, uuid.UUID, string, string) {
				return rt.UserID, rt.ClientID, rt.Token, rt.Scopes
			},
			isErr: true,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "refresh_tokens" WHERE user_id = $1 AND client_id = $2`)).
					WithArgs(rt.UserID, rt.ClientID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "refresh_tokens" ("token","user_id","client_id","scopes") VALUES ($1,$2,$3,$4)`)).
					WithArgs(rt.Token, rt.UserID, rt.ClientID, rt.Scopes).
					WillReturnError(testutil.ErrDB)
				mock.ExpectRollback()
			},
			arg: func() (uuid.UUID, uuid.UUID, string, string) {
				return rt.UserID, rt.ClientID, rt.Token, rt.Scopes
			},
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewRefreshTokenRepository(db)
			userID, clientID, token, scope := tt.arg()
			err := repo.Update(c, userID, clientID, token, scope)

			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
