package config

import (
	"strings"

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
	v.RegisterValidation("uuidStrings", isUUIDStringsValid)
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

func isUUIDStringsValid(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if field == "" {
		return true
	}

	UUIDs := strings.Split(field, ",")
	for _, ID := range UUIDs {
		if strings.TrimSpace(ID) == "" {
			continue
		}
		if _, err := uuid.Parse(strings.TrimSpace(ID)); err != nil {
			return false
		}
	}
	return true
}
