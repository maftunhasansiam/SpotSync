package handler

import "github.com/go-playground/validator/v10"

func formatValidationErrors(err error) map[string]string {
	errs := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errs[e.Field()] = e.Tag()
		}
	}
	return errs
}