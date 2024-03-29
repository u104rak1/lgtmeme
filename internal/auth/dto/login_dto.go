package dto

type LoginForm struct {
	Username string `form:"username" validate:"required,max=20"`
	Password string `form:"password" validate:"required,min=8,max=20"`
}
