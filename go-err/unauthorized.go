package go_err

import (
	"errors"
	"google.golang.org/grpc/codes"
	"net/http"
)

const (
	unauthorizedGRPCCode = codes.PermissionDenied
	unauthorizedHTTPCode = http.StatusUnauthorized
)

type UnauthorizedError struct {
	*BaseError
}

func NewUnauthorizedError(msg string) error {
	return &BaseError{
		Message:  msg,
		GRPCCode: unauthorizedGRPCCode,
		HTTPCode: unauthorizedHTTPCode,
	}
}

func NewUnauthorizedErrorWithCode(msg string, code int) error {
	return &BaseError{
		Code:     code,
		Message:  msg,
		GRPCCode: unauthorizedGRPCCode,
		HTTPCode: unauthorizedHTTPCode,
	}
}

func NewUnauthorizedErrorWithName(msg string, name string) error {
	return &BaseError{
		Name:     name,
		Message:  msg,
		GRPCCode: unauthorizedGRPCCode,
		HTTPCode: unauthorizedHTTPCode,
	}
}

func NewUnauthorizedErrorWithCodeAndName(msg string, code int, name string) error {
	return &BaseError{
		Code:     code,
		Name:     name,
		Message:  msg,
		GRPCCode: unauthorizedGRPCCode,
		HTTPCode: unauthorizedHTTPCode,
	}
}

func IsUnauthorizedErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == unauthorizedGRPCCode
}

func IsUnauthorizedError(err error) bool {
	if IsUnauthenticatedErrorGRPC(err) {
		return true
	}
	var expectedErr UnauthorizedError
	return errors.As(err, &expectedErr)
}
