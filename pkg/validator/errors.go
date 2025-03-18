package validator

import (
	"backoffice/internal/constants"
	"sync"
)

func (v *Validator) generateErrorMessage() {
	errorMessagesOnce.Do(func() {
		errorMessages = map[string]string{
			"required":         "field %s is required",
			"email_custom":     "email %s is not valid",
			"str_gt":           "field %s must have greater than %s characters",
			"str_lt":           "field %s must have less than %s characters",
			"has_lowercase":    "field %s must have at least one small character",
			"has_uppercase":    "field %s must have at least one big character",
			"has_special":      "field %s must have at least one special character",
			"oneof":            "field %s must have value one of allowed list: %s",
			"gte":              "field %s must be greater or equal than %s",
			"gt":               "field %s must be greater than %s",
			"url":              "field %s must be an url",
			"min":              "field %s must have at least %s characters",
			CustomDateTimeRule: "field %s must have datetime format: " + constants.TimeLayout,
			"eqfield":          "field %s confirmation does not match",
		}
	})
}

var (
	errorMessages     map[string]string
	errorMessagesOnce sync.Once
)
