package dto

type LoginForm struct {
	Username     string `form:"username" validate:"required,max=20"`
	Password     string `form:"password" validate:"required,min=8,max=20"`
	ScopeConsent string `form:"scopeConsent" validate:"oneof=true false"`
}

type LoginResp struct {
	RedirectURL string `json:"redirectURL"`
}
