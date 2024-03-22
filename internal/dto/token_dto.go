package dto

import "github.com/ucho456job/my_authn_authz/internal/model"

type TokenForm struct {
	GrantType    string         `form:"grant_type" validate:"required"`
	Code         string         `form:"code" validate:"required"`
	RedirectURI  string         `form:"redirect_uri" validate:"required"`
	ClientID     model.ClientID `form:"client_id" validate:"required"`
	ClientSecret string         `form:"client_secret" validate:"required"`
}
