package validator

import (
	"backoffice/internal/constants"
	"github.com/go-playground/validator/v10"
	"time"
)

const (
	CustomDateTimeRule = "custom_datetime"
)

func IsCustomDateTime(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}

	if _, err := time.Parse(constants.TimeLayout, fl.Field().String()); err != nil {
		return false
	}

	return true
}
