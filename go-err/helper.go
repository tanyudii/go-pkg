package go_err

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

const (
	metaKeyErrorCode = "Code"
	metaKeyErrorName = "name"
)

func FromStatus(s *status.Status) CustomError {
	if s == nil {
		return nil
	}

	base := &BaseError{
		GRPCCode: s.Code(),
		HTTPCode: HTTPStatusFromCode(s.Code()),
		Message:  s.Message(),
	}

	for _, detail := range s.Details() {

		//get error Code & name
		errInfo, valid := detail.(*errdetails.ErrorInfo)
		if valid {
			if codeStr, ok := errInfo.Metadata[metaKeyErrorCode]; ok {
				if code, err := strconv.Atoi(codeStr); err == nil {
					base.Code = code
				}
			}
			if name, ok := errInfo.Metadata[metaKeyErrorName]; ok {
				base.Name = name
			}
			continue
		}

		//get bad request details
		badRequest, valid := detail.(*errdetails.BadRequest)
		if valid && len(badRequest.FieldViolations) != 0 {
			base.Fields = make(ErrorField)
			for _, field := range badRequest.FieldViolations {
				base.Fields[field.Field] = field.Description
			}
			continue
		}
	}

	switch s.Code() {
	case badRequestGRPCCode:
		return &BadRequestError{BaseError: base}
	case notFoundGRPCCode:
		return &NotFoundError{BaseError: base}
	case tooManyRequestGRPCCode:
		return &TooManyRequestError{BaseError: base}
	case unauthenticatedGRPCCode:
		return &UnauthenticatedError{BaseError: base}
	case unauthorizedGRPCCode:
		return &UnauthorizedError{BaseError: base}
	default:
		return &InternalServerError{BaseError: base}
	}
}

func FromResponseError(r *ResponseError) error {
	if r == nil ||
		r.Status != ResponseErrorStatus ||
		r.Meta == nil {
		return nil
	}

	base := &BaseError{
		Code:     r.Meta.Code,
		Name:     r.Meta.Name,
		Message:  r.Message,
		GRPCCode: r.Meta.GrpcCode,
		HTTPCode: r.Meta.HttpCode,
		Fields:   r.Fields,
	}

	switch r.Meta.GrpcCode {
	case badRequestGRPCCode:
		return &BadRequestError{BaseError: base}
	case notFoundGRPCCode:
		return &NotFoundError{BaseError: base}
	case tooManyRequestGRPCCode:
		return &TooManyRequestError{BaseError: base}
	case unauthenticatedGRPCCode:
		return &UnauthenticatedError{BaseError: base}
	case unauthorizedGRPCCode:
		return &UnauthorizedError{BaseError: base}
	default:
		return &InternalServerError{BaseError: base}
	}
}

func GetErrorDetailsFromErrorGRPC(err error) []any {
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

func GetErrorCode(err error) int {
	details := GetErrorDetailsFromErrorGRPC(err)
	if len(details) == 0 {
		return 0
	}
	for _, detail := range details {
		errInfo, valid := detail.(*errdetails.ErrorInfo)
		if !valid {
			continue
		}
		if codeStr, ok := errInfo.Metadata[metaKeyErrorCode]; ok {
			if code, err := strconv.Atoi(codeStr); err == nil {
				return code
			}
		}
	}
	return 0
}

func GetErrorName(err error) string {
	details := GetErrorDetailsFromErrorGRPC(err)
	if len(details) == 0 {
		return ""
	}
	for _, detail := range details {
		errInfo, valid := detail.(*errdetails.ErrorInfo)
		if !valid {
			continue
		}
		if name, ok := errInfo.Metadata[metaKeyErrorName]; ok {
			return name
		}
	}
	return ""
}

func IsErrorCode(err error, code int) bool {
	return GetErrorCode(err) == code
}

func IsErrorName(err error, name string) bool {
	return GetErrorName(err) == name
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
