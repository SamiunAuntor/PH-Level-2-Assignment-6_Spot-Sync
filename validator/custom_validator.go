package validator

import (
	govalidator "github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *govalidator.Validate
}

func New() *CustomValidator {
	return &CustomValidator{
		validator: govalidator.New(),
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
