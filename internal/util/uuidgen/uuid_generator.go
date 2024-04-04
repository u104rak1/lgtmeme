package uuidgen

import "github.com/google/uuid"

type UUIDGenerator interface {
	New() uuid.UUID
}

type RealUUIDGenerator struct{}

func (g *RealUUIDGenerator) New() uuid.UUID {
	return uuid.New()
}

type MockUUIDGenerator struct {
	MockUUID uuid.UUID
}

func (g *MockUUIDGenerator) New() uuid.UUID {
	return g.MockUUID
}
