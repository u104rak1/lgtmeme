package repository_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/dto"
	"github.com/ucho456job/lgtmeme/internal/model"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/util/timer"
	"github.com/ucho456job/lgtmeme/test/testutil"
)

func TestImageRepository(t *testing.T) {
	gormDB, mock := testutil.SetupMockDB(t)
	mockTimer := timer.MockTimer{}

	testImages := []model.Image{
		{
			ID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			URL:       "http://example.com/image.jpg",
			Keyword:   "keyword1",
			UsedCount: 0,
			Reported:  false,
			Confirmed: false,
			CreatedAt: mockTimer.Now(),
		},
		{
			ID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			URL:       "http://example.com/image2.jpg",
			Keyword:   "keyword2",
			UsedCount: 1,
			Reported:  true,
			Confirmed: false,
			CreatedAt: mockTimer.Now(),
		},
		{
			ID:        uuid.MustParse("323e4567-e89b-12d3-a456-426614174000"),
			URL:       "http://example.com/image3.jpg",
			Keyword:   "keyword3",
			UsedCount: 2,
			Reported:  true,
			Confirmed: true,
			CreatedAt: mockTimer.Now(),
		},
	}
	i1 := testImages[0]
	i2 := testImages[1]
	i3 := testImages[2]

	sqlStatement := `INSERT INTO "images" ("url","keyword","used_count","reported","confirmed","id","created_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id","created_at"`

	createTests := []struct {
		name      string
		setupMock func()
		isErr     bool
	}{
		{
			name: "positive: Return nil, Create was successful",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(i1.URL, i1.Keyword, i1.UsedCount, i1.Reported, i1.Confirmed, i1.ID, i1.CreatedAt).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(i1.ID, i1.CreatedAt))
				mock.ExpectCommit()
			},
			isErr: false,
		},
		{
			name: "negative: Return error, because database connection error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(i1.URL, i1.Keyword, i1.UsedCount, i1.Reported, i1.Confirmed, i1.ID, i1.CreatedAt).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()
			},
			isErr: true,
		},
	}

	for _, tt := range createTests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(gormDB, &mockTimer)
			err := repo.Create(c, i1.ID, i1.URL, i1.Keyword)

			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	findImagesTests := []struct {
		name      string
		setupMock func()
		query     dto.GetImagesQuery
		result    *[]model.Image
		isErr     bool
	}{
		{
			name: "positive: Return images, with basic query",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "url", "keyword", "used_count", "reported", "confirmed", "created_at"}).
					AddRow(i1.ID, i1.URL, i1.Keyword, i1.UsedCount, i1.Reported, i1.Confirmed, i1.CreatedAt).
					AddRow(i3.ID, i3.URL, i3.Keyword, i3.UsedCount, i3.Reported, i3.Confirmed, i3.CreatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE confirmed = $1 OR reported = $2 ORDER BY created_at DESC LIMIT $3`)).
					WithArgs(true, false, config.GET_IMAGES_LIMIT).
					WillReturnRows(rows)
			},
			query: dto.GetImagesQuery{
				Page:             0,
				Keyword:          "",
				Sort:             "latest",
				FavoriteImageIDs: "",
				AuthCheck:        false,
			},
			result: &[]model.Image{i1, i3},
			isErr:  false,
		},
		{
			name: "positive: Return images, with favoriteImageIDs",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "url", "keyword", "used_count", "reported", "confirmed", "created_at"}).
					AddRow(i1.ID, i1.URL, i1.Keyword, i1.UsedCount, i1.Reported, i1.Confirmed, i1.CreatedAt).
					AddRow(i2.ID, i2.URL, i2.Keyword, i2.UsedCount, i2.Reported, i2.Confirmed, i2.CreatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE id IN ($1,$2) AND (confirmed = $3 OR reported = $4) ORDER BY created_at DESC LIMIT $5`)).
					WithArgs(i1.ID, i2.ID, true, false, config.GET_IMAGES_LIMIT).
					WillReturnRows(rows)
			},
			query: dto.GetImagesQuery{
				Page:             0,
				Keyword:          "",
				Sort:             "latest",
				FavoriteImageIDs: fmt.Sprintf("%s,%s", i1.ID, i2.ID),
				AuthCheck:        false,
			},
			result: &[]model.Image{i1, i2},
			isErr:  false,
		},
		{
			name: "positive: Return images, with keyword",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "url", "keyword", "used_count", "reported", "confirmed", "created_at"}).
					AddRow(i1.ID, i1.URL, i1.Keyword, i1.UsedCount, i1.Reported, i1.Confirmed, i1.CreatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE keyword LIKE $1 AND (confirmed = $2 OR reported = $3) ORDER BY created_at DESC LIMIT $4`)).
					WithArgs("%keyword1%", true, false, config.GET_IMAGES_LIMIT).
					WillReturnRows(rows)
			},
			query: dto.GetImagesQuery{
				Page:             0,
				Keyword:          "keyword1",
				Sort:             "latest",
				FavoriteImageIDs: "",
				AuthCheck:        false,
			},
			result: &[]model.Image{i1},
			isErr:  false,
		},
		{
			name: "positive: Return images, with authCheck",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "url", "keyword", "used_count", "reported", "confirmed", "created_at"}).
					AddRow(i2.ID, i2.URL, i2.Keyword, i2.UsedCount, i2.Reported, i2.Confirmed, i2.CreatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE confirmed = $1 AND reported = $2 ORDER BY created_at DESC LIMIT $3`)).
					WithArgs(false, true, config.GET_IMAGES_LIMIT).
					WillReturnRows(rows)
			},
			query: dto.GetImagesQuery{
				Page:             0,
				Keyword:          "",
				Sort:             "latest",
				FavoriteImageIDs: "",
				AuthCheck:        true,
			},
			result: &[]model.Image{i2},
			isErr:  false,
		},
		{
			name: "positive: Return images, with sort = popular",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "url", "keyword", "used_count", "reported", "confirmed", "created_at"}).
					AddRow(i3.ID, i3.URL, i3.Keyword, i3.UsedCount, i3.Reported, i3.Confirmed, i3.CreatedAt).
					AddRow(i2.ID, i2.URL, i2.Keyword, i2.UsedCount, i2.Reported, i2.Confirmed, i2.CreatedAt).
					AddRow(i1.ID, i1.URL, i1.Keyword, i1.UsedCount, i1.Reported, i1.Confirmed, i1.CreatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE confirmed = $1 OR reported = $2 ORDER BY used_count DESC, created_at DESC LIMIT $3`)).
					WithArgs(true, false, config.GET_IMAGES_LIMIT).
					WillReturnRows(rows)
			},
			query: dto.GetImagesQuery{
				Page:             0,
				Keyword:          "",
				Sort:             "popular",
				FavoriteImageIDs: "",
				AuthCheck:        false,
			},
			result: &[]model.Image{i3, i2, i1},
			isErr:  false,
		},
		{
			name: "negative: Return error, because database connection error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE confirmed = $1 OR reported = $2 ORDER BY created_at DESC LIMIT $3`)).
					WithArgs(true, false, config.GET_IMAGES_LIMIT).
					WillReturnError(errors.New("database connection failed"))
			},
			query: dto.GetImagesQuery{
				Page:             0,
				Keyword:          "",
				Sort:             "latest",
				FavoriteImageIDs: "",
				AuthCheck:        false,
			},
			result: nil,
			isErr:  true,
		},
	}

	for _, tt := range findImagesTests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(gormDB, &mockTimer)
			result, err := repo.FindImages(c, tt.query)

			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.result, result)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	findURLByIDTests := []struct {
		name      string
		setupMock func()
		result    bool
		isErr     bool
	}{
		{
			name: "positive: Return true, with existing ID",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "images" WHERE id = $1`)).
					WithArgs(i1.ID).
					WillReturnRows(rows)
			},
			result: true,
			isErr:  false,
		},
	}

	for _, tt := range findURLByIDTests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(gormDB, &mockTimer)
			result, err := repo.ExistsByID(c, i1.ID)

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
