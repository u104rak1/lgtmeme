package dto

import "github.com/google/uuid"

type TokenForm struct {
	GrantType    string    `form:"grant_type" validate:"required,grantType"`
	Code         string    `form:"code" validate:"required_with=GrantType=authorization_code"`
	RedirectURI  string    `form:"redirect_uri" validate:"required_with=GrantType=authorization_code"`
	ClientID     uuid.UUID `form:"client_id" validate:"required"`
	ClientSecret string    `form:"client_secret" validate:"required_with=GrantType=client_credentials|GrantType=refresh_token"`
}
