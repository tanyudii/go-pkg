package go_err

import "google.golang.org/grpc/codes"

const (
	ResponseErrorStatus = "error"
)

type ResponseError struct {
	Status  string     `json:"status,omitempty"`
	Message string     `json:"message,omitempty"`
	Meta    *ErrorMeta `json:"meta,omitempty"`
}

type ErrorMeta struct {
	Code     int        `json:"code,omitempty"`
	Name     string     `json:"name,omitempty"`
	GrpcCode codes.Code `json:"grpcCode,omitempty"`
	HttpCode int        `json:"httpCode,omitempty"`
	Fields   ErrorField `json:"fields,omitempty"`
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
			Fields:   c.GetFields(),
		},
	}
}

func NewResponseErrorWithFields(c CustomError, fields map[string]string) *ResponseError {
	resp := NewResponseError(c)
	resp.SetField(fields)
	return resp
}

func (r *ResponseError) SetField(fields map[string]string) {
	if r.Meta == nil {
		r.Meta = &ErrorMeta{}
	}
	r.Meta.Fields = fields
}

func (r *ResponseError) AddField(key, value string) {
	if r.Meta == nil {
		r.Meta = &ErrorMeta{}
	}
	if r.Meta.Fields == nil {
		r.Meta.Fields = make(ErrorField)
	}
	r.Meta.Fields[key] = value
}

func (r *ResponseError) ToBadRequest() *BadRequestError {
	return &BadRequestError{
		baseError: &baseError{
			code: r.Meta.Code,
		},
	}
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
