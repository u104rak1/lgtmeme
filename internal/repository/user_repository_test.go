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
	user = model.User{
		ID:       testutil.TestUUIDs[0],
		Name:     "name",
		Password: "password",
		Role:     "admin",
	}
)

func TestUserRepository_FirstByID(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	sqlStatement := `SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`

	tests := []struct {
		name      string
		setupMock func()
		arg       func() (uuid.UUID, []string)
		want      *model.User
		isErr     bool
	}{
		{
			name: "Return user, do not specify columns",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "password", "role"}).
					AddRow(user.ID, user.Name, user.Password, user.Role)
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(user.ID, 1).
					WillReturnRows(rows)
			},
			arg: func() (uuid.UUID, []string) {
				return user.ID, []string{}
			},
			want:  &user,
			isErr: false,
		},
		{
			name: "Return user, specify columns",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"role"}).
					AddRow(user.Role)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "role" FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(user.ID, 1).
					WillReturnRows(rows)
			},
			arg: func() (uuid.UUID, []string) {
				return user.ID, []string{"role"}
			},
			want:  &model.User{Role: user.Role},
			isErr: false,
		},
		{
			name: "Return error, because record not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(user.ID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			arg: func() (uuid.UUID, []string) {
				return user.ID, []string{}
			},
			want:  nil,
			isErr: true,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(user.ID, 1).
					WillReturnError(testutil.ErrDB)
			},
			arg: func() (uuid.UUID, []string) {
				return user.ID, []string{}
			},
			want:  nil,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewUserRepository(db)
			userID, columns := tt.arg()
			got, err := repo.FirstByID(c, userID, columns)

			if tt.isErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_FirstByName(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	sqlStatement := `SELECT * FROM "users" WHERE name = $1 ORDER BY "users"."id" LIMIT $2`

	tests := []struct {
		name      string
		setupMock func()
		arg       func() string
		want      *model.User
		isErr     bool
	}{
		{
			name: "Return user",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "password", "role"}).
					AddRow(user.ID, user.Name, user.Password, user.Role)
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(user.Name, 1).
					WillReturnRows(rows)
			},
			arg: func() string {
				return user.Name
			},
			want:  &user,
			isErr: false,
		},
		{
			name: "Return error, because record not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(user.Name, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			arg: func() string {
				return user.Name
			},
			want:  nil,
			isErr: true,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(user.Name, 1).
					WillReturnError(testutil.ErrDB)
			},
			arg: func() string {
				return user.Name
			},
			want:  nil,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewUserRepository(db)
			name := tt.arg()
			got, err := repo.FirstByName(c, name)

			if tt.isErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_ExistsByID(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	sqlStatement := `SELECT count(*) FROM "users" WHERE id = $1`

	tests := []struct {
		name      string
		setupMock func()
		arg       func() uuid.UUID
		want      bool
		isErr     bool
	}{
		{
			name: "Return true",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"count"}).
					AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(user.ID).
					WillReturnRows(rows)
			},
			arg: func() uuid.UUID {
				return user.ID
			},
			want:  true,
			isErr: false,
		},
		{
			name: "Return false",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"count"}).
					AddRow(0)
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(user.ID).
					WillReturnRows(rows)
			},
			arg: func() uuid.UUID {
				return user.ID
			},
			want:  false,
			isErr: false,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(user.ID).
					WillReturnError(testutil.ErrDB)
			},
			arg: func() uuid.UUID {
				return user.ID
			},
			want:  false,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewUserRepository(db)
			userID := tt.arg()
			got, err := repo.ExistsByID(c, userID)

			if tt.isErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
