package go_err

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type CustomError interface {
	Error() string
	GRPCStatus() *status.Status
	GetCode() int
	GetName() string
	GetGRPCCode() codes.Code
	GetHTTPCode() int
	GetFields() ErrorField
	SetFields(v ErrorField)
	GetData() any
	SetData(v any)
}

type BaseError struct {
	Code     int
	Name     string
	Message  string
	GRPCCode codes.Code
	HTTPCode int
	Data     any
	Fields   ErrorField
}

func (i *BaseError) Error() string {
	return i.Message
}

func (i *BaseError) GRPCStatus() *status.Status {
	stats := status.New(i.GRPCCode, i.Error())
	if customErr := i.getErrorInfoCustom(); customErr != nil {
		stats, _ = stats.WithDetails(customErr)
	}
	if fields := i.getBadRequestFields(); fields != nil {
		stats, _ = stats.WithDetails(fields)
	}
	return stats
}

func (i *BaseError) GetCode() int {
	return i.Code
}

func (i *BaseError) GetName() string {
	return i.Name
}

func (i *BaseError) GetGRPCCode() codes.Code {
	return i.GRPCCode
}

func (i *BaseError) GetHTTPCode() int {
	return i.HTTPCode
}

func (i *BaseError) GetFields() ErrorField {
	return i.Fields
}

func (i *BaseError) SetFields(v ErrorField) {
	i.Fields = v
}

func (i *BaseError) GetData() any {
	return i.Data
}

func (i *BaseError) SetData(v any) {
	i.Data = v
}

func (i *BaseError) getErrorInfoCustom() *errdetails.ErrorInfo {
	metaData := make(map[string]string)
	if i.Code != 0 {
		metaData[metaKeyErrorCode] = strconv.Itoa(i.Code)
	}
	if i.Name != "" {
		metaData[metaKeyErrorName] = i.Name
	}
	if len(metaData) == 0 {
		return nil
	}
	return &errdetails.ErrorInfo{Metadata: metaData}
}

func (i *BaseError) getBadRequestFields() *errdetails.BadRequest {
	var violations []*errdetails.BadRequest_FieldViolation
	for attr, msg := range i.Fields {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       attr,
			Description: msg,
		})
	}
	if len(violations) == 0 {
		return nil
	}
	return &errdetails.BadRequest{
		FieldViolations: violations,
	}
}
