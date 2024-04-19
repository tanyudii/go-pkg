package go_err

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

const (
	metaKeyErrorCode = "code"
	metaKeyErrorName = "name"
)

func FromStatus(s *status.Status) CustomError {
	if s == nil {
		return nil
	}

	base := &baseError{
		grpcCode: s.Code(),
		httpCode: HTTPStatusFromCode(s.Code()),
		message:  s.Message(),
	}

	for _, detail := range s.Details() {
		errInfo, valid := detail.(*errdetails.ErrorInfo)
		if valid {
			if codeStr, ok := errInfo.Metadata[metaKeyErrorCode]; ok {
				if code, err := strconv.Atoi(codeStr); err == nil {
					base.code = code
				}
			}
			if name, ok := errInfo.Metadata[metaKeyErrorName]; ok {
				base.name = name
			}
			continue
		}
		badRequest, valid := detail.(*errdetails.BadRequest)
		if valid {
			fields := make(ErrorField)
			for _, field := range badRequest.FieldViolations {
				fields[field.Field] = field.Description
			}
			return &BadRequestError{baseError: base, fields: fields}
		}
	}

	switch s.Code() {
	case badRequestGRPCCode:
		return &BadRequestError{baseError: base}
	case notFoundGRPCCode:
		return &NotFoundError{baseError: base}
	case tooManyRequestGRPCCode:
		return &TooManyRequestError{baseError: base}
	case unauthenticatedGRPCCode:
		return &UnauthenticatedError{baseError: base}
	case unauthorizedGRPCCode:
		return &UnauthorizedError{baseError: base}
	default:
		return &InternalServerError{baseError: base}
	}
}

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

func HTTPStatusFromCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		// Note, this deliberately doesn't translate to the similarly named '412 Precondition Failed' HTTP response status.
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
