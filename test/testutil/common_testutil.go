package testutil

import (
	"errors"

	"github.com/google/uuid"
)

var ErrDB = errors.New("database connection failed")

var TestUUIDs = []uuid.UUID{
	uuid.MustParse("00000000-0000-0000-0000-000000000000"),
	uuid.MustParse("11111111-1111-1111-1111-111111111111"),
	uuid.MustParse("22222222-2222-2222-2222-222222222222"),
	uuid.MustParse("33333333-3333-3333-3333-333333333333"),
	uuid.MustParse("44444444-4444-4444-4444-444444444444"),
	uuid.MustParse("55555555-5555-5555-5555-555555555555"),
	uuid.MustParse("66666666-6666-6666-6666-666666666666"),
	uuid.MustParse("77777777-7777-7777-7777-777777777777"),
	uuid.MustParse("88888888-8888-8888-8888-888888888888"),
	uuid.MustParse("99999999-9999-9999-9999-999999999999"),
}
