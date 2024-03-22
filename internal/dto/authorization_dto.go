package dto

import "github.com/google/uuid"

type AuthorizationQuery struct {
	ResponseType string    `validate:"required"`
	ClientID     uuid.UUID `validate:"required"`
	RedirectURI  string    `validate:"required,url"`
	Scope        string    `validate:"omitempty"`
	State        string    `validate:"omitempty"`
	Nonce        string    `validate:"omitempty"`
}
