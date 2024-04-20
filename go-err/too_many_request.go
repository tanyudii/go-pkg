package go_err

import (
	"errors"
	"google.golang.org/grpc/codes"
	"net/http"
)

const (
	tooManyRequestGRPCCode = codes.ResourceExhausted
	tooManyRequestHTTPCode = http.StatusTooManyRequests
)

type TooManyRequestError struct {
	*BaseError
}

func NewTooManyRequestError(msg string) error {
	return &TooManyRequestError{
		BaseError: &BaseError{
			Message:  msg,
			GRPCCode: tooManyRequestGRPCCode,
			HTTPCode: tooManyRequestHTTPCode,
		},
	}
}

func NewTooManyRequestErrorWithCode(msg string, code int) error {
	return &TooManyRequestError{
		BaseError: &BaseError{
			Code:     code,
			Message:  msg,
			GRPCCode: tooManyRequestGRPCCode,
			HTTPCode: tooManyRequestHTTPCode,
		},
	}
}

func NewTooManyRequestErrorWithName(msg string, name string) error {
	return &TooManyRequestError{
		BaseError: &BaseError{
			Name:     name,
			Message:  msg,
			GRPCCode: tooManyRequestGRPCCode,
			HTTPCode: tooManyRequestHTTPCode,
		},
	}
}

func NewTooManyRequestErrorWithCodeAndName(msg string, code int, name string) error {
	return &TooManyRequestError{
		BaseError: &BaseError{
			Code:     code,
			Name:     name,
			Message:  msg,
			GRPCCode: tooManyRequestGRPCCode,
			HTTPCode: tooManyRequestHTTPCode,
		},
	}
}

func IsTooManyRequestErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == tooManyRequestGRPCCode
}

func IsTooManyRequestError(err error) bool {
	if IsTooManyRequestErrorGRPC(err) {
		return true
	}
	var expectedErr *TooManyRequestError
	return errors.As(err, &expectedErr)
}
