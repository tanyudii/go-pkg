package go_err

import (
	"errors"
	"google.golang.org/grpc/codes"
	"net/http"
)

const (
	unauthenticatedGRPCCode = codes.Unauthenticated
	unauthenticatedHTTPCode = http.StatusUnauthorized
)

type UnauthenticatedError struct {
	*baseError
}

func NewUnauthenticatedError(msg string) error {
	return &UnauthenticatedError{
		baseError: &baseError{
			message:  msg,
			grpcCode: unauthenticatedGRPCCode,
			httpCode: unauthenticatedHTTPCode,
		},
	}
}

func NewUnauthenticatedErrorWithCode(msg string, code int) error {
	return &UnauthenticatedError{
		baseError: &baseError{
			code:     code,
			message:  msg,
			grpcCode: unauthenticatedGRPCCode,
			httpCode: unauthenticatedHTTPCode,
		},
	}
}

func NewUnauthenticatedErrorWithName(msg string, name string) error {
	return &UnauthenticatedError{
		baseError: &baseError{
			name:     name,
			message:  msg,
			grpcCode: unauthenticatedGRPCCode,
			httpCode: unauthenticatedHTTPCode,
		},
	}
}

func IsUnauthenticatedErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == unauthenticatedGRPCCode
}

func IsUnauthenticatedError(err error) bool {
	if IsUnauthenticatedErrorGRPC(err) {
		return true
	}
	var expectedErr *UnauthenticatedError
	return errors.As(err, &expectedErr)
}
