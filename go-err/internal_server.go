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
	*baseError
}

func NewInternalServerError(msg string) error {
	return &InternalServerError{
		baseError: &baseError{
			message:  msg,
			grpcCode: internalServerGRPCCode,
			httpCode: internalServerHTTPCode,
		},
	}
}

func NewInternalServerErrorWithCode(msg string, code int) error {
	return &InternalServerError{
		baseError: &baseError{
			code:     code,
			message:  msg,
			grpcCode: internalServerGRPCCode,
			httpCode: internalServerHTTPCode,
		},
	}
}

func NewInternalServerErrorWithName(msg string, name string) error {
	return &InternalServerError{
		baseError: &baseError{
			name:     name,
			message:  msg,
			grpcCode: internalServerGRPCCode,
			httpCode: internalServerHTTPCode,
		},
	}
}

func NewInternalServerErrorWithCodeAndName(msg string, code int, name string) error {
	return &InternalServerError{
		baseError: &baseError{
			code:     code,
			name:     name,
			message:  msg,
			grpcCode: internalServerGRPCCode,
			httpCode: internalServerHTTPCode,
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
