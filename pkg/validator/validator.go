package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func isForbiddenShortUrl(fl validator.FieldLevel) bool {
	return fl.Field().String() != "api"
}

func New() (*CustomValidator, error) {
	v := validator.New()

	err := v.RegisterValidation("short", isForbiddenShortUrl)
	if err != nil {
		return nil, fmt.Errorf("building custom short url validator: %w", err)
	}

	return &CustomValidator{v}, nil
}
