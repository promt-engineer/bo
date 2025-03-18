package rpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WrapInGRPCError(err error) error {
	validationErr, ok := err.(ValidationError)
	if ok {
		return status.Error(codes.InvalidArgument, validationErr.Error())
	}

	code, ok := errMap[err]
	if !ok {
		return status.Error(codes.Unknown, err.Error())
	}

	return status.Error(code, err.Error())
}
