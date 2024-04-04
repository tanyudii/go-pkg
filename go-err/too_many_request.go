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
	*baseError
}

func NewTooManyRequestError(msg string) error {
	return &TooManyRequestError{
		baseError: &baseError{
			message:  msg,
			grpcCode: tooManyRequestGRPCCode,
			httpCode: tooManyRequestHTTPCode,
		},
	}
}

func NewTooManyRequestErrorWithCode(msg string, code int) error {
	return &TooManyRequestError{
		baseError: &baseError{
			code:     code,
			message:  msg,
			grpcCode: tooManyRequestGRPCCode,
			httpCode: tooManyRequestHTTPCode,
		},
	}
}

func NewTooManyRequestErrorWithName(msg string, name string) error {
	return &TooManyRequestError{
		baseError: &baseError{
			name:     name,
			message:  msg,
			grpcCode: tooManyRequestGRPCCode,
			httpCode: tooManyRequestHTTPCode,
		},
	}
}

func NewTooManyRequestErrorWithCodeAndName(msg string, code int, name string) error {
	return &TooManyRequestError{
		baseError: &baseError{
			code:     code,
			name:     name,
			message:  msg,
			grpcCode: tooManyRequestGRPCCode,
			httpCode: tooManyRequestHTTPCode,
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
