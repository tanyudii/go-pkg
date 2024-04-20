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
	*BaseError
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{
		&BaseError{
			Message:  msg,
			GRPCCode: notFoundGRPCCode,
			HTTPCode: notFoundHTTPCode,
		},
	}
}

func NewNotFoundErrorWithCode(msg string, code int) error {
	return &NotFoundError{
		&BaseError{
			Code:     code,
			Message:  msg,
			GRPCCode: notFoundGRPCCode,
			HTTPCode: notFoundHTTPCode,
		},
	}
}

func NewNotFoundErrorWithName(msg string, name string) error {
	return &NotFoundError{
		&BaseError{
			Name:     name,
			Message:  msg,
			GRPCCode: notFoundGRPCCode,
			HTTPCode: notFoundHTTPCode,
		},
	}
}

func NewNotFoundErrorWithCodeAndName(msg string, code int, name string) error {
	return &NotFoundError{
		&BaseError{
			Code:     code,
			Name:     name,
			Message:  msg,
			GRPCCode: notFoundGRPCCode,
			HTTPCode: notFoundHTTPCode,
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
