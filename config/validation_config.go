package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewValidator() echo.Validator {
	v := validator.New()
	v.RegisterValidation("grantType", isGrantTypeValid)
	v.RegisterValidation("sort", isSortValid)
	v.RegisterValidation("uuidSlice", isUUIDSliceValid)
	return &CustomValidator{validator: v}
}

func isGrantTypeValid(fl validator.FieldLevel) bool {
	grantType := fl.Field().String()
	allowedGrantTypes := []string{"authorization_code", "client_credentials", "refresh_token"}
	for _, agt := range allowedGrantTypes {
		if grantType == agt {
			return true
		}
	}
	return false
}

func isSortValid(fl validator.FieldLevel) bool {
	sort := fl.Field().String()
	allowedSortTypes := []string{"latest", "popular"}
	for _, ast := range allowedSortTypes {
		if sort == ast {
			return true
		}
	}
	return false
}

func isUUIDSliceValid(fl validator.FieldLevel) bool {
	UUIDs := fl.Field().Interface().([]string)
	for _, ID := range UUIDs {
		if _, err := uuid.Parse(ID); err != nil {
			return false
		}
	}
	return true
}
