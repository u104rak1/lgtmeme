package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func InitValidator() echo.Validator {
	v := validator.New()
	v.RegisterValidation("grantType", isGrantTypeValid)
	return &CustomValidator{validator: v}
}

func isGrantTypeValid(fl validator.FieldLevel) bool {
	grantType := fl.Field().String()
	allowedGrantTypes := []string{"authorization_code", "client_credentials", "refresh_token"}
	for _, allowedGrantType := range allowedGrantTypes {
		if grantType == allowedGrantType {
			return true
		}
	}
	return false
}
