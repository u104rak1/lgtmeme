package dto

import "github.com/google/uuid"

type AuthzQuery struct {
	ResponseType string    `validate:"required"`
	ClientID     uuid.UUID `validate:"required"`
	RedirectURI  string    `validate:"required,url"`
	Scope        string    `validate:"omitempty"`
	State        string    `validate:"required"`
	Nonce        string    `validate:"omitempty"`
}
