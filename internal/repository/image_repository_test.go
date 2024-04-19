package repository_test

import (
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
	"gorm.io/gorm"
)

var (
	mockTimer  = timer.MockTimer{}
	testImages = []model.Image{
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
	i1 = testImages[0]
	i2 = testImages[1]
	i3 = testImages[2]
)

func TestImageRepository_Create(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	sqlStatement := `INSERT INTO "images" ("url","keyword","used_count","reported","confirmed","id","created_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id","created_at"`

	tests := []struct {
		name      string
		setupMock func()
		isErr     bool
	}{
		{
			name: "positive: Create was successful",
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
					WillReturnError(testutil.ErrDB)
				mock.ExpectRollback()
			},
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(db, &mockTimer)
			err := repo.Create(c, i1.ID, i1.URL, i1.Keyword)

			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestImageRepository_FindImages(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	tests := []struct {
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
			name: "positive: Return empty, with no images",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "url", "keyword", "used_count", "reported", "confirmed", "created_at"})
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
			result: &[]model.Image{},
			isErr:  false,
		},
		{
			name: "negative: Return error, because database connection error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE confirmed = $1 OR reported = $2 ORDER BY created_at DESC LIMIT $3`)).
					WithArgs(true, false, config.GET_IMAGES_LIMIT).
					WillReturnError(testutil.ErrDB)
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(db, &mockTimer)
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
}

func TestImageRepository_FindURLByID(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	tests := []struct {
		name      string
		setupMock func()
		id        uuid.UUID
		result    *string
		isErr     bool
	}{
		{
			name: "positive: Return url",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"url"}).AddRow(i1.URL)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "url" FROM "images" WHERE id = $1 ORDER BY "images"."id" LIMIT $2`)).
					WithArgs(i1.ID, 1).
					WillReturnRows(rows)
			},
			id:     i1.ID,
			result: &i1.URL,
			isErr:  false,
		},
		{
			name: "negative: Return nil, because record not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "url" FROM "images" WHERE id = $1 ORDER BY "images"."id" LIMIT $2`)).
					WithArgs(i1.ID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			id:     i1.ID,
			result: nil,
			isErr:  true,
		},
		{
			name: "negative: Return error, because database connection error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "url" FROM "images" WHERE id = $1 ORDER BY "images"."id" LIMIT $2`)).
					WithArgs(i1.ID, 1).
					WillReturnError(testutil.ErrDB)
			},
			id:     i1.ID,
			result: nil,
			isErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(db, &mockTimer)
			result, err := repo.FindURLByID(c, tt.id)

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

func TestImageRepository_ExistsByID(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	nonExistingID := uuid.MustParse("423e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name      string
		setupMock func()
		id        uuid.UUID
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
			id:     i1.ID,
			result: true,
			isErr:  false,
		},
		{
			name: "positive: Return false, with non-existing ID",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "images" WHERE id = $1`)).
					WithArgs(nonExistingID).
					WillReturnRows(rows)
			},
			id:     nonExistingID,
			result: false,
			isErr:  false,
		},
		{
			name: "negative: Return error, because database connection error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "images" WHERE id = $1`)).
					WithArgs(i1.ID).
					WillReturnError(testutil.ErrDB)
			},
			id:     i1.ID,
			result: false,
			isErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(db, &mockTimer)
			result, err := repo.ExistsByID(c, tt.id)

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

func TestImageRepository_Update(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	tests := []struct {
		name      string
		setupMock func()
		id        uuid.UUID
		reqType   dto.PatchImageReqType
		isErr     bool
	}{
		{
			name: "positive: Update was successful, with PatchImageReqTypeUsed",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "used_count"=used_count + $1 WHERE id = $2`)).
					WithArgs(1, i1.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			id:      i1.ID,
			reqType: dto.PatchImageReqTypeUsed,
			isErr:   false,
		},
		{
			name: "positive: Update was successful, with PatchImageReqTypeReport",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "reported"=$1 WHERE id = $2`)).
					WithArgs(true, i1.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			id:      i1.ID,
			reqType: dto.PatchImageReqTypeReport,
			isErr:   false,
		},
		{
			name: "positive: Update was successful, with PatchImageReqTypeConfirm",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "confirmed"=$1 WHERE id = $2`)).
					WithArgs(true, i2.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			id:      i2.ID,
			reqType: dto.PatchImageReqTypeConfirm,
			isErr:   false,
		},
		{
			name:      "negative: Return error, because PatchImageReqType is invalid",
			setupMock: func() {},
			id:        i1.ID,
			reqType:   "invalid",
			isErr:     true,
		},
		{
			name: "negative: Return error, because record not found",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "used_count"=used_count + $1 WHERE id = $2`)).
					WithArgs(1, i1.ID).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			id:      i1.ID,
			reqType: dto.PatchImageReqTypeUsed,
			isErr:   true,
		},
		{
			name: "negative: Return error, because database connection error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "used_count"=used_count + $1 WHERE id = $2`)).
					WithArgs(1, i1.ID).
					WillReturnError(testutil.ErrDB)
				mock.ExpectRollback()
			},
			id:      i1.ID,
			reqType: dto.PatchImageReqTypeUsed,
			isErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(db, &mockTimer)
			err := repo.Update(c, tt.id, tt.reqType)

			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestImageRepository_Delete(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	tests := []struct {
		name      string
		setupMock func()
		id        uuid.UUID
		isErr     bool
	}{
		{
			name: "positive: Delete was successful",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "images" WHERE id = $1`)).
					WithArgs(i1.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			id:    i1.ID,
			isErr: false,
		},
		{
			name: "negative: Return error, because record not found",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "images" WHERE id = $1`)).
					WithArgs(i1.ID).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			id:    i1.ID,
			isErr: true,
		},
		{
			name: "negative: Return error, because database connection error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "images" WHERE id = $1`)).
					WithArgs(i1.ID).
					WillReturnError(testutil.ErrDB)
				mock.ExpectRollback()
			},
			id:    i1.ID,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(db, &mockTimer)
			err := repo.Delete(c, tt.id)

			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
