package go_err

import (
	"google.golang.org/grpc/codes"
)

const (
	ResponseErrorStatus = "error"
)

type ErrorField map[string]string

func (f ErrorField) getFirstAndTotal() (string, int) {
	t := len(f)
	if t > 0 {
		t--
	}
	for k := range f {
		return f[k], t
	}
	return "", 0
}

type ResponseError struct {
	Status  string     `json:"status,omitempty"`
	Message string     `json:"message,omitempty"`
	Meta    *ErrorMeta `json:"meta,omitempty"`
	Fields  ErrorField `json:"fields,omitempty"`
}

type ErrorMeta struct {
	Code     int        `json:"Code,omitempty"`
	Name     string     `json:"name,omitempty"`
	GrpcCode codes.Code `json:"grpcCode,omitempty"`
	HttpCode int        `json:"httpCode,omitempty"`
}

func NewResponseError(c CustomError) *ResponseError {
	return &ResponseError{
		Status:  ResponseErrorStatus,
		Message: c.Error(),
		Meta: &ErrorMeta{
			Code:     c.GetCode(),
			Name:     c.GetName(),
			GrpcCode: c.GetGRPCCode(),
			HttpCode: c.GetHTTPCode(),
		},
		Fields: c.GetFields(),
	}
}

func (r *ResponseError) SetField(fields map[string]string) {
	r.Fields = fields
}

func (r *ResponseError) AddField(key, value string) {
	if r.Fields == nil {
		r.Fields = make(ErrorField)
	}
	r.Fields[key] = value
}
