package dto

type AuthorizationQuery struct {
	ResponseType string
	ClientID     string
	RedirectURI  string
	Scope        string
	State        string
	Nonce        string
}
