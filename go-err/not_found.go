package go_err

import (
	"errors"
	"google.golang.org/grpc/codes"
	"net/http"
)

const (
	notFoundGRPCCode = codes.NotFound
	notFoundHTTPCode = http.StatusNotFound
)

type NotFoundError struct {
	*baseError
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{
		baseError: &baseError{
			message:  msg,
			grpcCode: notFoundGRPCCode,
			httpCode: notFoundHTTPCode,
		},
	}
}

func NewNotFoundErrorWithCode(msg string, code int) error {
	return &NotFoundError{
		baseError: &baseError{
			code:     code,
			message:  msg,
			grpcCode: notFoundGRPCCode,
			httpCode: notFoundHTTPCode,
		},
	}
}

func NewNotFoundErrorWithName(msg string, name string) error {
	return &NotFoundError{
		baseError: &baseError{
			name:     name,
			message:  msg,
			grpcCode: notFoundGRPCCode,
			httpCode: notFoundHTTPCode,
		},
	}
}

func IsNotFoundErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == notFoundGRPCCode
}

func IsNotFoundError(err error) bool {
	if IsNotFoundErrorGRPC(err) {
		return true
	}
	var expectedErr *NotFoundError
	return errors.As(err, &expectedErr)
}
