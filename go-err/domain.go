package go_err

import "google.golang.org/grpc/codes"

type CustomError interface {
	Error() string
	GetCode() int
	GetName() string
	GetGRPCCode() codes.Code
	GetHTTPCode() int
	ToResponseError() *ResponseError
}
