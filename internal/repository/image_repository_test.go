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
			ID:        testutil.TestUUIDs[0],
			URL:       "http://example.com/image.jpg",
			Keyword:   "keyword1",
			UsedCount: 0,
			Reported:  false,
			Confirmed: false,
			CreatedAt: mockTimer.Now(),
		},
		{
			ID:        testutil.TestUUIDs[1],
			URL:       "http://example.com/image2.jpg",
			Keyword:   "keyword2",
			UsedCount: 1,
			Reported:  true,
			Confirmed: false,
			CreatedAt: mockTimer.Now(),
		},
		{
			ID:        testutil.TestUUIDs[2],
			URL:       "http://example.com/image3.jpg",
			Keyword:   "keyword3",
			UsedCount: 2,
			Reported:  true,
			Confirmed: true,
			CreatedAt: mockTimer.Now(),
		},
	}
	img1 = testImages[0]
	img2 = testImages[1]
	img3 = testImages[2]
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
			name: "Create was successful",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(img1.URL, img1.Keyword, img1.UsedCount, img1.Reported, img1.Confirmed, img1.ID, img1.CreatedAt).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(img1.ID, img1.CreatedAt))
				mock.ExpectCommit()
			},
			isErr: false,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(img1.URL, img1.Keyword, img1.UsedCount, img1.Reported, img1.Confirmed, img1.ID, img1.CreatedAt).
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
			err := repo.Create(c, img1.ID, img1.URL, img1.Keyword)

			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestImageRepository_FindByGetImagesQuery(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	tests := []struct {
		name      string
		setupMock func()
		arg       func() dto.GetImagesQuery
		want      *[]model.Image
		isErr     bool
	}{
		{
			name: "Return images, with basic query",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "url", "keyword", "used_count", "reported", "confirmed", "created_at"}).
					AddRow(img1.ID, img1.URL, img1.Keyword, img1.UsedCount, img1.Reported, img1.Confirmed, img1.CreatedAt).
					AddRow(img3.ID, img3.URL, img3.Keyword, img3.UsedCount, img3.Reported, img3.Confirmed, img3.CreatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, url FROM "images" WHERE confirmed = $1 OR reported = $2 ORDER BY created_at DESC LIMIT $3`)).
					WithArgs(true, false, config.GET_IMAGES_LIMIT).
					WillReturnRows(rows)
			},
			arg: func() dto.GetImagesQuery {
				return dto.GetImagesQuery{
					Page:             0,
					Keyword:          "",
					Sort:             "latest",
					FavoriteImageIDs: "",
					AuthCheck:        false,
				}
			},
			want:  &[]model.Image{img1, img3},
			isErr: false,
		},
		{
			name: "Return images, with favoriteImageIDs",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "url", "keyword", "used_count", "reported", "confirmed", "created_at"}).
					AddRow(img1.ID, img1.URL, img1.Keyword, img1.UsedCount, img1.Reported, img1.Confirmed, img1.CreatedAt).
					AddRow(img2.ID, img2.URL, img2.Keyword, img2.UsedCount, img2.Reported, img2.Confirmed, img2.CreatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, url FROM "images" WHERE id IN ($1,$2) AND (confirmed = $3 OR reported = $4) ORDER BY created_at DESC LIMIT $5`)).
					WithArgs(img1.ID, img2.ID, true, false, config.GET_IMAGES_LIMIT).
					WillReturnRows(rows)
			},
			arg: func() dto.GetImagesQuery {
				return dto.GetImagesQuery{
					Page:             0,
					Keyword:          "",
					Sort:             "latest",
					FavoriteImageIDs: fmt.Sprintf("%s,%s", img1.ID, img2.ID),
					AuthCheck:        false,
				}
			},
			want:  &[]model.Image{img1, img2},
			isErr: false,
		},
		{
			name: "Return images, with keyword",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "url", "keyword", "used_count", "reported", "confirmed", "created_at"}).
					AddRow(img1.ID, img1.URL, img1.Keyword, img1.UsedCount, img1.Reported, img1.Confirmed, img1.CreatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, url FROM "images" WHERE keyword LIKE $1 AND (confirmed = $2 OR reported = $3) ORDER BY created_at DESC LIMIT $4`)).
					WithArgs("%keyword1%", true, false, config.GET_IMAGES_LIMIT).
					WillReturnRows(rows)
			},
			arg: func() dto.GetImagesQuery {
				return dto.GetImagesQuery{
					Page:             0,
					Keyword:          "keyword1",
					Sort:             "latest",
					FavoriteImageIDs: "",
					AuthCheck:        false,
				}
			},
			want:  &[]model.Image{img1},
			isErr: false,
		},
		{
			name: "Return images, with authCheck",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "url", "keyword", "used_count", "reported", "confirmed", "created_at"}).
					AddRow(img2.ID, img2.URL, img2.Keyword, img2.UsedCount, img2.Reported, img2.Confirmed, img2.CreatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, url FROM "images" WHERE confirmed = $1 AND reported = $2 ORDER BY created_at DESC LIMIT $3`)).
					WithArgs(false, true, config.GET_IMAGES_LIMIT).
					WillReturnRows(rows)
			},
			arg: func() dto.GetImagesQuery {
				return dto.GetImagesQuery{
					Page:             0,
					Keyword:          "",
					Sort:             "latest",
					FavoriteImageIDs: "",
					AuthCheck:        true,
				}
			},
			want:  &[]model.Image{img2},
			isErr: false,
		},
		{
			name: "Return images, with sort is popular",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "url", "keyword", "used_count", "reported", "confirmed", "created_at"}).
					AddRow(img3.ID, img3.URL, img3.Keyword, img3.UsedCount, img3.Reported, img3.Confirmed, img3.CreatedAt).
					AddRow(img2.ID, img2.URL, img2.Keyword, img2.UsedCount, img2.Reported, img2.Confirmed, img2.CreatedAt).
					AddRow(img1.ID, img1.URL, img1.Keyword, img1.UsedCount, img1.Reported, img1.Confirmed, img1.CreatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, url FROM "images" WHERE confirmed = $1 OR reported = $2 ORDER BY used_count DESC, created_at DESC LIMIT $3`)).
					WithArgs(true, false, config.GET_IMAGES_LIMIT).
					WillReturnRows(rows)
			},
			arg: func() dto.GetImagesQuery {
				return dto.GetImagesQuery{
					Page:             0,
					Keyword:          "",
					Sort:             "popular",
					FavoriteImageIDs: "",
					AuthCheck:        false,
				}
			},
			want:  &[]model.Image{img3, img2, img1},
			isErr: false,
		},
		{
			name: "Return empty",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "url", "keyword", "used_count", "reported", "confirmed", "created_at"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, url FROM "images" WHERE confirmed = $1 OR reported = $2 ORDER BY created_at DESC LIMIT $3`)).
					WithArgs(true, false, config.GET_IMAGES_LIMIT).
					WillReturnRows(rows)
			},
			arg: func() dto.GetImagesQuery {
				return dto.GetImagesQuery{
					Page:             0,
					Keyword:          "",
					Sort:             "latest",
					FavoriteImageIDs: "",
					AuthCheck:        false,
				}
			},
			want:  &[]model.Image{},
			isErr: false,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, url FROM "images" WHERE confirmed = $1 OR reported = $2 ORDER BY created_at DESC LIMIT $3`)).
					WithArgs(true, false, config.GET_IMAGES_LIMIT).
					WillReturnError(testutil.ErrDB)
			},
			arg: func() dto.GetImagesQuery {
				return dto.GetImagesQuery{
					Page:             0,
					Keyword:          "",
					Sort:             "latest",
					FavoriteImageIDs: "",
					AuthCheck:        false,
				}
			},
			want:  nil,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(db, &mockTimer)
			q := tt.arg()
			actual, err := repo.FindByGetImagesQuery(c, q)

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

func TestImageRepository_FirstByID(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	sqlStatement := `SELECT * FROM "images" WHERE id = $1 ORDER BY "images"."id" LIMIT $2`

	tests := []struct {
		name      string
		setupMock func()
		arg       func() (uuid.UUID, []string)
		want      *model.Image
		isErr     bool
	}{
		{
			name: "Return image, do not specify columns",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "url", "keyword", "used_count", "reported", "confirmed", "created_at"}).
					AddRow(img1.ID, img1.URL, img1.Keyword, img1.UsedCount, img1.Reported, img1.Confirmed, img1.CreatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(img1.ID, 1).
					WillReturnRows(rows)
			},
			arg: func() (uuid.UUID, []string) {
				return img1.ID, []string{}
			},
			want:  &img1,
			isErr: false,
		},
		{
			name: "Return image, specify columns",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"url"}).AddRow(img1.URL)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "url" FROM "images" WHERE id = $1 ORDER BY "images"."id" LIMIT $2`)).
					WithArgs(img1.ID, 1).
					WillReturnRows(rows)
			},
			arg: func() (uuid.UUID, []string) {
				return img1.ID, []string{"url"}
			},
			want:  &model.Image{URL: img1.URL},
			isErr: false,
		},
		{
			name: "Return error, because record not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(img1.ID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			arg: func() (uuid.UUID, []string) {
				return img1.ID, []string{}
			},
			want:  nil,
			isErr: true,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(sqlStatement)).
					WithArgs(img1.ID, 1).
					WillReturnError(testutil.ErrDB)
			},
			arg: func() (uuid.UUID, []string) {
				return img1.ID, []string{}
			},
			want:  nil,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(db, &mockTimer)
			id, columns := tt.arg()
			actual, err := repo.FirstByID(c, id, columns)

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

func TestImageRepository_ExistsByID(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	nonExistingID := testutil.TestUUIDs[3]

	tests := []struct {
		name      string
		setupMock func()
		arg       func() uuid.UUID
		want      bool
		isErr     bool
	}{
		{
			name: "Return true, with existing ID",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "images" WHERE id = $1`)).
					WithArgs(img1.ID).
					WillReturnRows(rows)
			},
			arg: func() uuid.UUID {
				return img1.ID
			},
			want:  true,
			isErr: false,
		},
		{
			name: "Return false, with non-existing ID",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "images" WHERE id = $1`)).
					WithArgs(nonExistingID).
					WillReturnRows(rows)
			},
			arg: func() uuid.UUID {
				return nonExistingID
			},
			want:  false,
			isErr: false,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "images" WHERE id = $1`)).
					WithArgs(img1.ID).
					WillReturnError(testutil.ErrDB)
			},
			arg: func() uuid.UUID {
				return img1.ID
			},
			want:  false,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(db, &mockTimer)
			id := tt.arg()
			actual, err := repo.ExistsByID(c, id)

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

func TestImageRepository_Update(t *testing.T) {
	db, mock := testutil.SetupMockDB(t)

	tests := []struct {
		name      string
		setupMock func()
		arg       func() (uuid.UUID, dto.PatchImageReqType)
		isErr     bool
	}{
		{
			name: "Update was successful, with PatchImageReqTypeUsed",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "used_count"=used_count + $1 WHERE id = $2`)).
					WithArgs(1, img1.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			arg: func() (uuid.UUID, dto.PatchImageReqType) {
				return img1.ID, dto.PatchImageReqTypeUsed
			},
			isErr: false,
		},
		{
			name: "Update was successful, with PatchImageReqTypeReport",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "reported"=$1 WHERE id = $2`)).
					WithArgs(true, img1.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			arg: func() (uuid.UUID, dto.PatchImageReqType) {
				return img1.ID, dto.PatchImageReqTypeReport
			},
			isErr: false,
		},
		{
			name: "Update was successful, with PatchImageReqTypeConfirm",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "confirmed"=$1 WHERE id = $2`)).
					WithArgs(true, img2.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			arg: func() (uuid.UUID, dto.PatchImageReqType) {
				return img2.ID, dto.PatchImageReqTypeConfirm
			},
			isErr: false,
		},
		{
			name:      "Return error, because PatchImageReqType is invalid",
			setupMock: func() {},
			arg: func() (uuid.UUID, dto.PatchImageReqType) {
				return img1.ID, "invalid"
			},
			isErr: true,
		},
		{
			name: "Return error, because record not found",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "used_count"=used_count + $1 WHERE id = $2`)).
					WithArgs(1, img1.ID).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			arg: func() (uuid.UUID, dto.PatchImageReqType) {
				return img1.ID, dto.PatchImageReqTypeUsed
			},
			isErr: true,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "images" SET "used_count"=used_count + $1 WHERE id = $2`)).
					WithArgs(1, img1.ID).
					WillReturnError(testutil.ErrDB)
				mock.ExpectRollback()
			},
			arg: func() (uuid.UUID, dto.PatchImageReqType) {
				return img1.ID, dto.PatchImageReqTypeUsed
			},
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(db, &mockTimer)
			id, reqType := tt.arg()
			err := repo.Update(c, id, reqType)

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
		arg       func() uuid.UUID
		isErr     bool
	}{
		{
			name: "Delete was successful",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "images" WHERE id = $1`)).
					WithArgs(img1.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			arg: func() uuid.UUID {
				return img1.ID
			},
			isErr: false,
		},
		{
			name: "Return error, because record not found",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "images" WHERE id = $1`)).
					WithArgs(img1.ID).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			arg: func() uuid.UUID {
				return img1.ID
			},
			isErr: true,
		},
		{
			name: "Return error, because db error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "images" WHERE id = $1`)).
					WithArgs(img1.ID).
					WillReturnError(testutil.ErrDB)
				mock.ExpectRollback()
			},
			arg: func() uuid.UUID {
				return img1.ID
			},
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testutil.SetupMinEchoContext()
			tt.setupMock()
			repo := repository.NewImageRepository(db, &mockTimer)
			id := tt.arg()
			err := repo.Delete(c, id)

			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
