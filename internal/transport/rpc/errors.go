package rpc

import (
	e "backoffice/internal/errors"
	"backoffice/pkg/validator"
	"bytes"
	"google.golang.org/grpc/codes"
)

var errMap = map[error]codes.Code{
	e.ErrNotAuthorized:         codes.Unauthenticated,
	e.ErrDoesNotHavePermission: codes.PermissionDenied,
}

type ValidationError struct {
	ErrorsBug []validator.TaggedError
}

func NewValidationError(errFromValidator error) ValidationError {
	return ValidationError{
		ErrorsBug: validator.CheckValidationErrors(errFromValidator),
	}
}

func (v ValidationError) Error() string {
	buf := bytes.Buffer{}

	for i, err := range v.ErrorsBug {
		buf.WriteString(err.Err.Error())

		if i != len(v.ErrorsBug)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}
