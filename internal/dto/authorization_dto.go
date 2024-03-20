package dto

type AuthoraizationQuery struct {
	ResponseType string
	ClientID     string
	RedirectURI  string
	Scope        string
	State        string
	Nonce        string
}
