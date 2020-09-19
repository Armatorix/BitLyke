package validator

import (
	"fmt"
	"net/url"

	"gopkg.in/go-playground/validator.v9"
)

func isURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

type CustomValidator struct {
	Validator *validator.Validate
}

func isURLString(fl validator.FieldLevel) bool {
	return isURL(fl.Field().String())
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func New() (*CustomValidator, error) {
	v := validator.New()
	err := v.RegisterValidation("url", isURLString)
	if err != nil {
		return nil, fmt.Errorf("building custom validator: %w", err)
	}
	return &CustomValidator{v}, nil
}
