package go_err

import (
	"errors"
	"net/http"
	"testing"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewBadRequestError(t *testing.T) {
	msg := "bad request"

	err := NewBadRequestError(msg)

	if err == nil {
		t.Fatal("NewBadRequestError() returned nil")
	}

	// Type assert to access CustomError methods
	customErr, ok := err.(CustomError)
	if !ok {
		t.Fatal("Error does not implement CustomError interface")
	}

	if customErr.GetGRPCCode() != codes.InvalidArgument {
		t.Errorf("GetGRPCCode() = %v; want InvalidArgument", customErr.GetGRPCCode())
	}

	if customErr.GetHTTPCode() != http.StatusBadRequest {
		t.Errorf("GetHTTPCode() = %d; want %d", customErr.GetHTTPCode(), http.StatusBadRequest)
	}

	if customErr.Error() != msg {
		t.Errorf("Error() = %s; want %s", customErr.Error(), msg)
	}
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("not found")

	if err == nil {
		t.Fatal("NewNotFoundError() returned nil")
	}

	customErr, ok := err.(CustomError)
	if !ok {
		t.Fatal("Error does not implement CustomError interface")
	}

	if customErr.GetGRPCCode() != codes.NotFound {
		t.Errorf("GetGRPCCode() = %v; want NotFound", customErr.GetGRPCCode())
	}

	if customErr.GetHTTPCode() != http.StatusNotFound {
		t.Errorf("GetHTTPCode() = %d; want %d", customErr.GetHTTPCode(), http.StatusNotFound)
	}
}

func TestNewNotFoundErrorWithCodeAndName(t *testing.T) {
	code := 404
	name := "RESOURCE_NOT_FOUND"
	msg := "Resource not found"

	err := NewNotFoundErrorWithCodeAndName(msg, code, name)

	if err == nil {
		t.Fatal("NewNotFoundErrorWithCodeAndName() returned nil")
	}

	customErr, ok := err.(CustomError)
	if !ok {
		t.Fatal("Error does not implement CustomError interface")
	}

	if customErr.GetCode() != code {
		t.Errorf("GetCode() = %d; want %d", customErr.GetCode(), code)
	}

	if customErr.GetName() != name {
		t.Errorf("GetName() = %s; want %s", customErr.GetName(), name)
	}
}

func TestNewInternalServerError(t *testing.T) {
	msg := "Internal server error"

	err := NewInternalServerError(msg)

	if err == nil {
		t.Fatal("NewInternalServerError() returned nil")
	}

	customErr, ok := err.(CustomError)
	if !ok {
		t.Fatal("Error does not implement CustomError interface")
	}

	if customErr.GetGRPCCode() != codes.Internal {
		t.Errorf("GetGRPCCode() = %v; want Internal", customErr.GetGRPCCode())
	}

	if customErr.GetHTTPCode() != http.StatusInternalServerError {
		t.Errorf("GetHTTPCode() = %d; want %d", customErr.GetHTTPCode(), http.StatusInternalServerError)
	}
}

func TestNewInternalServerErrorWithCodeAndName(t *testing.T) {
	code := 500
	name := "INTERNAL_ERROR"
	msg := "Something went wrong"

	err := NewInternalServerErrorWithCodeAndName(msg, code, name)

	if err == nil {
		t.Fatal("NewInternalServerErrorWithCodeAndName() returned nil")
	}

	customErr, ok := err.(CustomError)
	if !ok {
		t.Fatal("Error does not implement CustomError interface")
	}

	if customErr.GetCode() != code {
		t.Errorf("GetCode() = %d; want %d", customErr.GetCode(), code)
	}

	if customErr.GetName() != name {
		t.Errorf("GetName() = %s; want %s", customErr.GetName(), name)
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	err := NewUnauthorizedError("unauthorized access")

	if err == nil {
		t.Fatal("NewUnauthorizedError() returned nil")
	}

	customErr, ok := err.(CustomError)
	if !ok {
		t.Fatal("Error does not implement CustomError interface")
	}

	if customErr.GetHTTPCode() != http.StatusUnauthorized {
		t.Errorf("GetHTTPCode() = %d; want %d", customErr.GetHTTPCode(), http.StatusUnauthorized)
	}
}

func TestNewUnauthenticatedError(t *testing.T) {
	err := NewUnauthenticatedError("authentication required")

	if err == nil {
		t.Fatal("NewUnauthenticatedError() returned nil")
	}

	customErr, ok := err.(CustomError)
	if !ok {
		t.Fatal("Error does not implement CustomError interface")
	}

	if customErr.GetGRPCCode() != codes.Unauthenticated {
		t.Errorf("GetGRPCCode() = %v; want Unauthenticated", customErr.GetGRPCCode())
	}

	if customErr.GetHTTPCode() != http.StatusUnauthorized {
		t.Errorf("GetHTTPCode() = %d; want %d", customErr.GetHTTPCode(), http.StatusUnauthorized)
	}
}

func TestNewTooManyRequestError(t *testing.T) {
	err := NewTooManyRequestError("rate limit exceeded")

	if err == nil {
		t.Fatal("NewTooManyRequestError() returned nil")
	}

	customErr, ok := err.(CustomError)
	if !ok {
		t.Fatal("Error does not implement CustomError interface")
	}

	if customErr.GetGRPCCode() != codes.ResourceExhausted {
		t.Errorf("GetGRPCCode() = %v; want ResourceExhausted", customErr.GetGRPCCode())
	}

	if customErr.GetHTTPCode() != http.StatusTooManyRequests {
		t.Errorf("GetHTTPCode() = %d; want %d", customErr.GetHTTPCode(), http.StatusTooManyRequests)
	}
}

func TestBadRequestErrorWithFields(t *testing.T) {
	fields := ErrorField{
		"email":    "Invalid email format",
		"password": "Password too short",
	}
	err := NewBadRequestErrorUsingFieldsOrNil(fields)

	if err == nil {
		t.Fatal("NewBadRequestErrorUsingFieldsOrNil() returned nil")
	}

	customErr, ok := err.(CustomError)
	if !ok {
		t.Fatal("Error does not implement CustomError interface")
	}

	resultFields := customErr.GetFields()
	if len(resultFields) != 2 {
		t.Errorf("GetFields() length = %d; want 2", len(resultFields))
	}

	if resultFields["email"] != "Invalid email format" {
		t.Errorf("GetFields()[\"email\"] = %s; want 'Invalid email format'", resultFields["email"])
	}

	if resultFields["password"] != "Password too short" {
		t.Errorf("GetFields()[\"password\"] = %s; want 'Password too short'", resultFields["password"])
	}
}

func TestErrorSetFields(t *testing.T) {
	fields := ErrorField{
		"username": "Username is required",
	}
	err := NewBadRequestErrorWithFields("validation failed", fields)

	customErr, ok := err.(CustomError)
	if !ok {
		t.Fatal("Error does not implement CustomError interface")
	}

	// Now test setting new fields
	newFields := ErrorField{
		"password": "Password too short",
	}

	customErr.SetFields(newFields)

	resultFields := customErr.GetFields()
	if len(resultFields) != 1 {
		t.Errorf("GetFields() length = %d; want 1", len(resultFields))
	}

	if resultFields["password"] != "Password too short" {
		t.Errorf("GetFields()[\"password\"] = %s; want 'Password too short'", resultFields["password"])
	}
}

func TestIsNotFoundError(t *testing.T) {
	err := NewNotFoundError("resource not found")

	if !IsNotFoundError(err) {
		t.Error("IsNotFoundError() should return true for NotFoundError")
	}

	otherErr := NewInternalServerError("internal error")
	if IsNotFoundError(otherErr) {
		t.Error("IsNotFoundError() should return false for InternalServerError")
	}

	// Test with plain error
	plainErr := errors.New("not found")
	if IsNotFoundError(plainErr) {
		t.Error("IsNotFoundError() should return false for plain error")
	}
}

func TestIsInternalServerError(t *testing.T) {
	err := NewInternalServerError("internal error")

	if !IsInternalServerError(err) {
		t.Error("IsInternalServerError() should return true for InternalServerError")
	}

	otherErr := NewNotFoundError("not found")
	if IsInternalServerError(otherErr) {
		t.Error("IsInternalServerError() should return false for NotFoundError")
	}
}

func TestGRPCStatusConversion(t *testing.T) {
	err := NewNotFoundError("resource not found")

	customErr, ok := err.(CustomError)
	if !ok {
		t.Fatal("Error does not implement CustomError interface")
	}

	grpcStatus := customErr.GRPCStatus()

	if grpcStatus == nil {
		t.Fatal("GRPCStatus() returned nil")
	}

	if grpcStatus.Code() != codes.NotFound {
		t.Errorf("GRPCStatus().Code() = %v; want NotFound", grpcStatus.Code())
	}

	if grpcStatus.Message() == "" {
		t.Error("GRPCStatus().Message() should not be empty")
	}
}

func TestIsErrorCode(t *testing.T) {
	err := NewNotFoundErrorWithCodeAndName("not found", 404, "NOT_FOUND")

	if !IsErrorCode(err, 404) {
		t.Error("IsErrorCode() should return true for matching error code")
	}

	if IsErrorCode(err, 500) {
		t.Error("IsErrorCode() should return false for non-matching error code")
	}
}

func TestIsErrorName(t *testing.T) {
	err := NewNotFoundErrorWithName("not found", "RESOURCE_NOT_FOUND")

	if !IsErrorName(err, "RESOURCE_NOT_FOUND") {
		t.Error("IsErrorName() should return true for matching error name")
	}

	if IsErrorName(err, "OTHER_ERROR") {
		t.Error("IsErrorName() should return false for non-matching error name")
	}
}

func TestHTTPStatusFromCode(t *testing.T) {
	tests := []struct {
		code         codes.Code
		expectedHTTP int
	}{
		{codes.OK, http.StatusOK},
		{codes.NotFound, http.StatusNotFound},
		{codes.InvalidArgument, http.StatusBadRequest},
		{codes.Unauthenticated, http.StatusUnauthorized},
		{codes.PermissionDenied, http.StatusForbidden},
		{codes.ResourceExhausted, http.StatusTooManyRequests},
		{codes.Internal, http.StatusInternalServerError},
		{codes.Unavailable, http.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.code.String(), func(t *testing.T) {
			result := HTTPStatusFromCode(tt.code)
			if result != tt.expectedHTTP {
				t.Errorf("HTTPStatusFromCode(%v) = %d; want %d", tt.code, result, tt.expectedHTTP)
			}
		})
	}
}

func TestFromStatus(t *testing.T) {
	// Create a gRPC status with error info
	st, _ := status.New(codes.InvalidArgument, "validation failed").
		WithDetails(&errdetails.ErrorInfo{
			Metadata: map[string]string{
				"code": "400",
				"name": "VALIDATION_ERROR",
			},
		})

	customErr := FromStatus(st)

	if customErr == nil {
		t.Fatal("FromStatus() returned nil")
	}

	if customErr.GetHTTPCode() != http.StatusBadRequest {
		t.Errorf("GetHTTPCode() = %d; want %d", customErr.GetHTTPCode(), http.StatusBadRequest)
	}

	if customErr.GetName() != "VALIDATION_ERROR" {
		t.Errorf("GetName() = %s; want 'VALIDATION_ERROR'", customErr.GetName())
	}

	if customErr.GetCode() != 400 {
		t.Errorf("GetCode() = %d; want 400", customErr.GetCode())
	}
}

func TestGetErrorDetailsFromErrorGRPC(t *testing.T) {
	err := NewNotFoundError("not found")

	details := GetErrorDetailsFromErrorGRPC(err)

	if details == nil {
		t.Error("GetErrorDetailsFromErrorGRPC() returned nil")
	}
}

func TestGetErrorGRPCCodeFromErrorGRPC(t *testing.T) {
	err := NewNotFoundError("not found")

	code := GetErrorGRPCCodeFromErrorGRPC(err)

	if code != codes.NotFound {
		t.Errorf("GetErrorGRPCCodeFromErrorGRPC() = %v; want NotFound", code)
	}
}

func TestCustomErrorInterface(t *testing.T) {
	// Test that BaseError implements CustomError
	baseErr := &BaseError{
		Code:     400,
		Name:     "TEST_ERROR",
		Message:  "Test error message",
		GRPCCode: codes.InvalidArgument,
		HTTPCode: http.StatusBadRequest,
		Fields:   nil,
	}

	var _ CustomError = baseErr

	if baseErr.Error() != "Test error message" {
		t.Errorf("Error() = %s; want 'Test error message'", baseErr.Error())
	}

	if baseErr.GetCode() != 400 {
		t.Errorf("GetCode() = %d; want 400", baseErr.GetCode())
	}

	if baseErr.GetName() != "TEST_ERROR" {
		t.Errorf("GetName() = %s; want 'TEST_ERROR'", baseErr.GetName())
	}
}

func TestErrorField(t *testing.T) {
	fields := ErrorField{
		"field1": "error1",
		"field2": "error2",
	}

	if len(fields) != 2 {
		t.Errorf("ErrorField length = %d; want 2", len(fields))
	}

	// Test that we can modify the field
	fields["field3"] = "error3"

	if len(fields) != 3 {
		t.Errorf("Modified ErrorField length = %d; want 3", len(fields))
	}
}
