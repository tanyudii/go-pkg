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
	*BaseError
}

func NewBadRequestError(msg string) error {
	return &BadRequestError{
		&BaseError{
			Message:  msg,
			GRPCCode: badRequestGRPCCode,
			HTTPCode: badRequestHTTPCode,
		},
	}
}

func NewBadRequestErrorWithCode(msg string, code int) error {
	return &BadRequestError{
		&BaseError{
			Code:     code,
			Message:  msg,
			GRPCCode: badRequestGRPCCode,
			HTTPCode: badRequestHTTPCode,
		},
	}
}

func NewBadRequestErrorWithName(msg string, name string) error {
	return &BadRequestError{
		&BaseError{
			Name:     name,
			Message:  msg,
			GRPCCode: badRequestGRPCCode,
			HTTPCode: badRequestHTTPCode,
		},
	}
}

func NewBadRequestErrorWithCodeAndName(msg string, code int, name string) error {
	return &BadRequestError{
		&BaseError{
			Code:     code,
			Name:     name,
			Message:  msg,
			GRPCCode: badRequestGRPCCode,
			HTTPCode: badRequestHTTPCode,
		},
	}
}

func NewBadRequestErrorWithFields(msg string, fields ErrorField) error {
	return &BadRequestError{
		&BaseError{
			Message:  msg,
			GRPCCode: badRequestGRPCCode,
			HTTPCode: badRequestHTTPCode,
			Fields:   fields,
		},
	}
}

func NewBadRequestErrorWithCodeAndFields(msg string, code int, fields ErrorField) error {
	return &BadRequestError{
		&BaseError{
			Code:     code,
			Message:  msg,
			GRPCCode: badRequestGRPCCode,
			HTTPCode: badRequestHTTPCode,
			Fields:   fields,
		},
	}
}

func NewBadRequestErrorWithNameAndFields(msg string, name string, fields ErrorField) error {
	return &BadRequestError{
		&BaseError{
			Name:     name,
			Message:  msg,
			GRPCCode: badRequestGRPCCode,
			HTTPCode: badRequestHTTPCode,
			Fields:   fields,
		},
	}
}

func NewBadRequestWithCodeNameAndFields(msg string, code int, name string, fields ErrorField) error {
	return &BadRequestError{
		&BaseError{
			Code:     code,
			Name:     name,
			Message:  msg,
			GRPCCode: badRequestGRPCCode,
			HTTPCode: badRequestHTTPCode,
			Fields:   fields,
		},
	}
}

func NewBadRequestErrorUsingFieldsOrNil(fields ErrorField) error {
	if len(fields) == 0 {
		return nil
	}
	first, total := fields.getFirstAndTotal()
	if total >= 1 {
		return NewBadRequestErrorWithFields(fmt.Sprintf("%s. and there are %d errors", first, total), fields)
	}
	return NewBadRequestErrorWithFields(first, fields)
}

func IsBadRequestErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == badRequestGRPCCode
}

func IsBadRequestError(err error) bool {
	if IsBadRequestErrorGRPC(err) {
		return true
	}
	var expectedErr BadRequestError
	return errors.As(err, &expectedErr)
}
