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
	*baseError
}

func NewUnauthorizedError(msg string) error {
	return &UnauthorizedError{
		baseError: &baseError{
			message:  msg,
			grpcCode: unauthorizedGRPCCode,
			httpCode: unauthorizedHTTPCode,
		},
	}
}

func NewUnauthorizedErrorWithCode(msg string, code int) error {
	return &UnauthorizedError{
		baseError: &baseError{
			code:     code,
			message:  msg,
			grpcCode: unauthorizedGRPCCode,
			httpCode: unauthorizedHTTPCode,
		},
	}
}

func NewUnauthorizedErrorWithName(msg string, name string) error {
	return &UnauthorizedError{
		baseError: &baseError{
			name:     name,
			message:  msg,
			grpcCode: unauthorizedGRPCCode,
			httpCode: unauthorizedHTTPCode,
		},
	}
}

func IsUnauthorizedErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == unauthorizedGRPCCode
}

func IsUnauthorizedError(err error) bool {
	if IsUnauthenticatedErrorGRPC(err) {
		return true
	}
	var expectedErr *UnauthorizedError
	return errors.As(err, &expectedErr)
}
