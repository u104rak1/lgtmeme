package dto

type AuthorizationQuery struct {
	ResponseType string `validate:"required"`
	ClientID     string `validate:"required"`
	RedirectURI  string `validate:"required,url"`
	Scope        string `validate:"omitempty"`
	State        string `validate:"omitempty"`
	Nonce        string `validate:"omitempty"`
}
