package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "building custom short url validator")
	}

	return &CustomValidator{v}, nil
}
