package go_grpc

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	goerr "pkg.tanyudii.me/go-pkg/go-err"
	gologger "pkg.tanyudii.me/go-pkg/go-logger"
	"strings"
)

const (
	HeaderContentType          = "Content-Type"
	HeaderAccept               = "Accept"
	HeaderAcceptLanguage       = "accept-language"
	HeaderAuthorization        = "authorization"
	HeaderInternalCallPassword = "internalcallpassword"
	HeaderKeyUserID            = "userid"
	HeaderKeyUserType          = "usertype"
	HeaderKeyCompanyID         = "companyid"
	HeaderKeyClientID          = "clientid"
	HeaderKeyRequestID         = "requestid"
	HeaderUserAgent            = "user-agent"
	HeaderGRPCUserAgent        = "grpcgateway-user-agent"
)

var (
	mapHeaderTransform = map[string]string{
		HeaderContentType:          HeaderContentType,
		HeaderAccept:               HeaderAccept,
		HeaderAcceptLanguage:       HeaderAcceptLanguage,
		HeaderInternalCallPassword: HeaderInternalCallPassword,
		HeaderKeyUserID:            HeaderKeyUserID,
		HeaderKeyUserType:          HeaderKeyUserType,
		HeaderKeyCompanyID:         HeaderKeyCompanyID,
		HeaderKeyClientID:          HeaderKeyClientID,
		HeaderKeyRequestID:         HeaderKeyRequestID,
		HeaderUserAgent:            HeaderGRPCUserAgent,
	}
)

func MuxCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{HeaderContentType, HeaderAccept, HeaderAuthorization}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	gologger.Infof("preflight request for %s\n", r.URL.Path)
}

func MuxHandleRoutingError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, httpStatus int) {
	if httpStatus != http.StatusMethodNotAllowed {
		runtime.DefaultRoutingErrorHandler(ctx, mux, marshaler, w, r, httpStatus)
		return
	}
	err := &runtime.HTTPStatusError{
		HTTPStatus: httpStatus,
		Err:        status.Error(codes.Unimplemented, http.StatusText(httpStatus)),
	}
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}

func MuxErrorHandler(_ context.Context, _ *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	s := status.Convert(err)

	customStatus := goerr.FromStatus(s)

	var resp *goerr.ResponseError
	if customStatus != nil {
		resp = goerr.NewResponseError(customStatus)
	} else {
		var custom goerr.CustomError
		errors.As(goerr.NewInternalServerError("internal error"), &custom)
		resp = goerr.NewResponseError(custom)
	}

	w.Header().Del("Trailer")
	w.Header().Del("Transfer-Encoding")

	contentType := m.ContentType(resp)
	w.Header().Set("Content-Type", contentType)

	if s.Code() == codes.Unauthenticated {
		w.Header().Set("WWW-Authenticate", s.Message())
	}

	buf, merr := m.Marshal(resp)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", s, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = io.WriteString(w, `{"status": "Internal Server Error", "message": "failed to marshal error message"}`); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	st := runtime.HTTPStatusFromCode(s.Code())
	w.WriteHeader(st)
	if _, err = w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
}

func MuxIncomingHeaderMatcher(key string) (string, bool) {
	if h, ok := mapHeaderTransform[strings.ToLower(key)]; ok {
		return h, true
	}
	return "", false
}

func MuxHandleRoutingRedirect(_ context.Context, w http.ResponseWriter, _ proto.Message) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])
		w.WriteHeader(http.StatusFound)
	}
	return nil
}
