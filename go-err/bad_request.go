package go_err

import (
	"errors"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

const (
	badRequestGRPCCode = codes.InvalidArgument
	badRequestHTTPCode = http.StatusBadRequest
)

type BadRequestError struct {
	*baseError
	fields ErrorField
}

type ErrorField map[string]string

func (f ErrorField) GetFirstErrorAndOtherTotal() (string, int) {
	total := len(f)
	if total > 0 {
		total--
	}
	for k := range f {
		return f[k], total
	}
	return "", 0
}

func (i *BadRequestError) GRPCStatus() *status.Status {
	stats := status.New(i.GetGRPCCode(), i.Error())
	if fields := i.GetBadRequestFields(); fields != nil {
		stats, _ = stats.WithDetails(fields)
	}
	if customErr := i.GetErrorInfoCustom(); customErr != nil {
		stats, _ = stats.WithDetails(customErr)
	}
	return stats
}

func (i *BadRequestError) GetFields() ErrorField {
	return i.fields
}

func (i *BadRequestError) GetBadRequestFields() *errdetails.BadRequest {
	errFields := i.GetFields()
	if len(errFields) == 0 {
		return nil
	}
	br := &errdetails.BadRequest{}
	for attr, msg := range i.GetFields() {
		br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       attr,
			Description: msg,
		})
	}
	return br
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
		fields: fields,
		baseError: &baseError{
			message:  msg,
			grpcCode: badRequestGRPCCode,
			httpCode: badRequestHTTPCode,
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
