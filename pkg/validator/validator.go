package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Validate(s interface{}) []ValidationError {
	var errs []ValidationError

	if err := validate.Struct(s); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, ValidationError{
				Field:   e.Field(),
				Message: msgForTag(e.Tag(), e.Param()),
			})
		}
	}

	return errs
}

func msgForTag(tag string, param string) string {
	switch tag {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email format"
	case "min":
		return "value is too short, minimum " + param
	case "max":
		return "value is too long, maximum " + param
	case "oneof":
		return "value must be one of: " + param
	case "url":
		return "invalid url format"
	}
	return "invalid value"
}
