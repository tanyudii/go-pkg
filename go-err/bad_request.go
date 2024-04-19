package go_err

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"net/http"
)

const (
	badRequestGRPCCode = codes.InvalidArgument
	badRequestHTTPCode = http.StatusBadRequest
)

type BadRequestError struct {
	*baseError
}

func NewBadRequestError(msg string) error {
	return &BadRequestError{
		baseError: &baseError{
			message:  msg,
			grpcCode: badRequestGRPCCode,
			httpCode: badRequestHTTPCode,
		},
	}
}

func NewBadRequestErrorWithCode(msg string, code int) error {
	return &BadRequestError{
		baseError: &baseError{
			code:     code,
			message:  msg,
			grpcCode: badRequestGRPCCode,
			httpCode: badRequestHTTPCode,
		},
	}
}

func NewBadRequestErrorWithName(msg string, name string) error {
	return &BadRequestError{
		baseError: &baseError{
			name:     name,
			message:  msg,
			grpcCode: badRequestGRPCCode,
			httpCode: badRequestHTTPCode,
		},
	}
}

func NewBadRequestErrorWithCodeAndName(msg string, code int, name string) error {
	return &BadRequestError{
		baseError: &baseError{
			code:     code,
			name:     name,
			message:  msg,
			grpcCode: badRequestGRPCCode,
			httpCode: badRequestHTTPCode,
		},
	}
}

func NewBadRequestErrorWithFields(msg string, fields ErrorField) error {
	return &BadRequestError{
		baseError: &baseError{
			message:  msg,
			grpcCode: badRequestGRPCCode,
			httpCode: badRequestHTTPCode,
			fields:   fields,
		},
	}
}

func NewBadRequestErrorUsingFieldsOrNil(fields ErrorField) error {
	if len(fields) == 0 {
		return nil
	}
	firstErr, otherErr := fields.GetFirstErrorAndOtherTotal()
	if otherErr >= 1 {
		return NewBadRequestErrorWithFields(fmt.Sprintf("%s. and there are %d errors", firstErr, otherErr), fields)
	}
	return NewBadRequestErrorWithFields(firstErr, fields)
}

func IsBadRequestErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == badRequestGRPCCode
}

func IsBadRequestError(err error) bool {
	if IsBadRequestErrorGRPC(err) {
		return true
	}
	var expectedErr *BadRequestError
	return errors.As(err, &expectedErr)
}
