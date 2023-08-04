package go_err

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type CustomError interface {
	Error() string
	GetCode() int
	GetName() string
	GetGRPCCode() codes.Code
	GetHTTPCode() int
}

type baseError struct {
	code     int
	name     string
	message  string
	grpcCode codes.Code
	httpCode int
}

func (i *baseError) Error() string {
	return i.message
}

func (i *baseError) GetCode() int {
	return i.code
}

func (i *baseError) GetName() string {
	return i.name
}

func (i *baseError) GetGRPCCode() codes.Code {
	return i.grpcCode
}

func (i *baseError) GetHTTPCode() int {
	return i.httpCode
}

func (i *baseError) GRPCStatus() *status.Status {
	stats := status.New(i.GetGRPCCode(), i.Error())
	if customErr := i.GetErrorInfoCustom(); customErr != nil {
		stats, _ = stats.WithDetails(customErr)
	}
	return stats
}

func (i *baseError) GetErrorInfoCustom() *errdetails.ErrorInfo {
	metaData := make(map[string]string)
	if code := i.GetCode(); code != 0 {
		metaData[metaKeyErrorCode] = strconv.Itoa(code)
	}
	if name := i.GetName(); name != "" {
		metaData[metaKeyErrorName] = name
	}
	if len(metaData) == 0 {
		return nil
	}
	return &errdetails.ErrorInfo{
		Metadata: metaData,
	}
}
