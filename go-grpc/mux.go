package go_grpc

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"net/http"
	gologger "pkg.tanyudii.me/go-pkg/go-logger"
	"strings"
)

const (
	HeaderContentType   = "Content-Type"
	HeaderAccept        = "Accept"
	HeaderAuthorization = "authorization"
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

func MuxErrorHandler(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
	s := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(s.Code())
	newError := runtime.HTTPStatusError{HTTPStatus: httpStatus, Err: err}
	runtime.DefaultHTTPErrorHandler(ctx, mux, m, w, req, &newError)
}

func MuxHandleRoutingRedirect(_ context.Context, w http.ResponseWriter, _ proto.Message) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])
		w.WriteHeader(http.StatusFound)
	}
	return nil
}
