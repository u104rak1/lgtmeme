package dto

type TokenForm struct {
	GrantType    string `form:"grant_type" validate:"required"`
	Code         string `form:"code" validate:"required"`
	RedirectURI  string `form:"redirect_uri" validate:"required"`
	ClientID     string `form:"client_id" validate:"required"`
	ClientSecret string `form:"client_secret" validate:"required"`
}
