package go_err

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

const (
	metaKeyErrorCode = "code"
	metaKeyErrorName = "name"
)

func CustomErrorToMapInterface(err CustomError) map[string]interface{} {
	obj := make(map[string]interface{})
	if code := err.GetCode(); code != 0 {
		obj[metaKeyErrorCode] = code
	}
	if name := err.GetName(); name != "" {
		obj[metaKeyErrorName] = name
	}
	return obj
}

func GetErrorDetailsFromErrorGRPC(err error) []interface{} {
	e, ok := status.FromError(err)
	if !ok {
		return nil
	}
	return e.Details()
}

func GetErrorGRPCCodeFromErrorGRPC(err error) codes.Code {
	e, ok := status.FromError(err)
	if !ok {
		return codes.Unknown
	}
	return e.Code()
}

func IsErrorCode(err error, code int) bool {
	details := GetErrorDetailsFromErrorGRPC(err)
	if len(details) == 0 {
		return false
	}
	codeStr := strconv.Itoa(code)
	for _, detail := range details {
		errInfo, valid := detail.(*errdetails.ErrorInfo)
		if !valid {
			continue
		}
		if errInfo.Metadata[metaKeyErrorCode] == codeStr {
			return true
		}
	}
	return false
}

func IsErrorName(err error, name string) bool {
	details := GetErrorDetailsFromErrorGRPC(err)
	if len(details) == 0 {
		return false
	}
	for _, detail := range details {
		errInfo, valid := detail.(*errdetails.ErrorInfo)
		if !valid {
			continue
		}
		if errInfo.Metadata[metaKeyErrorName] == name {
			return true
		}
	}
	return false
}
