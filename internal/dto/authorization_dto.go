package dto

import "github.com/ucho456job/my_authn_authz/internal/model"

type AuthorizationQuery struct {
	ResponseType string         `validate:"required"`
	ClientID     model.ClientID `validate:"required"`
	RedirectURI  string         `validate:"required,url"`
	Scope        string         `validate:"omitempty"`
	State        string         `validate:"omitempty"`
	Nonce        string         `validate:"omitempty"`
}
