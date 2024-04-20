package go_err

import (
	"errors"
	"google.golang.org/grpc/codes"
	"net/http"
)

const (
	internalServerGRPCCode = codes.Internal
	internalServerHTTPCode = http.StatusInternalServerError
)

type InternalServerError struct {
	*BaseError
}

func NewInternalServerError(msg string) error {
	return &InternalServerError{
		&BaseError{
			Message:  msg,
			GRPCCode: internalServerGRPCCode,
			HTTPCode: internalServerHTTPCode,
		},
	}
}

func NewInternalServerErrorWithCode(msg string, code int) error {
	return &InternalServerError{
		&BaseError{
			Code:     code,
			Message:  msg,
			GRPCCode: internalServerGRPCCode,
			HTTPCode: internalServerHTTPCode,
		},
	}
}

func NewInternalServerErrorWithName(msg string, name string) error {
	return &InternalServerError{
		&BaseError{
			Name:     name,
			Message:  msg,
			GRPCCode: internalServerGRPCCode,
			HTTPCode: internalServerHTTPCode,
		},
	}
}

func NewInternalServerErrorWithCodeAndName(msg string, code int, name string) error {
	return &InternalServerError{
		&BaseError{
			Code:     code,
			Name:     name,
			Message:  msg,
			GRPCCode: internalServerGRPCCode,
			HTTPCode: internalServerHTTPCode,
		},
	}
}

func IsInternalServerErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == internalServerGRPCCode
}

func IsInternalServerError(err error) bool {
	if IsInternalServerErrorGRPC(err) {
		return true
	}
	var expectedErr *InternalServerError
	return errors.As(err, &expectedErr)
}
