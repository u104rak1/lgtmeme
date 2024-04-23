package testutil

import (
	"errors"

	"github.com/google/uuid"
)

var ErrDB = errors.New("database connection failed")

var TestUUIDs = []uuid.UUID{
	uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
	uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
	uuid.MustParse("323e4567-e89b-12d3-a456-426614174000"),
	uuid.MustParse("423e4567-e89b-12d3-a456-426614174000"),
	uuid.MustParse("523e4567-e89b-12d3-a456-426614174000"),
	uuid.MustParse("623e4567-e89b-12d3-a456-426614174000"),
	uuid.MustParse("723e4567-e89b-12d3-a456-426614174000"),
	uuid.MustParse("823e4567-e89b-12d3-a456-426614174000"),
	uuid.MustParse("923e4567-e89b-12d3-a456-426614174000"),
	uuid.MustParse("023e4567-e89b-12d3-a456-426614174000"),
}
