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
	*BaseError
}

func NewUnauthenticatedError(msg string) error {
	return &UnauthenticatedError{
		&BaseError{
			Message:  msg,
			GRPCCode: unauthenticatedGRPCCode,
			HTTPCode: unauthenticatedHTTPCode,
		},
	}
}

func NewUnauthenticatedErrorWithCode(msg string, code int) error {
	return &UnauthenticatedError{
		&BaseError{
			Code:     code,
			Message:  msg,
			GRPCCode: unauthenticatedGRPCCode,
			HTTPCode: unauthenticatedHTTPCode,
		},
	}
}

func NewUnauthenticatedErrorWithName(msg string, name string) error {
	return &UnauthenticatedError{
		&BaseError{
			Name:     name,
			Message:  msg,
			GRPCCode: unauthenticatedGRPCCode,
			HTTPCode: unauthenticatedHTTPCode,
		},
	}
}

func NewUnauthenticatedErrorWithCodeAndName(msg string, code int, name string) error {
	return &UnauthenticatedError{
		&BaseError{
			Code:     code,
			Name:     name,
			Message:  msg,
			GRPCCode: unauthenticatedGRPCCode,
			HTTPCode: unauthenticatedHTTPCode,
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
