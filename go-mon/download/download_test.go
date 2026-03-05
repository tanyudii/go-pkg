package download

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		serverStatus   int
		method         string
		want           string
		wantErr        bool
		errContains    string
	}{
		{
			name:           "successful GET request",
			serverResponse: "Hello, World!",
			serverStatus:   200,
			method:         http.MethodGet,
			want:           "Hello, World!",
			wantErr:        false,
		},
		{
			name:           "successful POST request",
			serverResponse: `{"status":"ok"}`,
			serverStatus:   201,
			method:         http.MethodPost,
			want:           `{"status":"ok"}`,
			wantErr:        false,
		},
		{
			name:           "successful PUT request",
			serverResponse: "updated",
			serverStatus:   200,
			method:         http.MethodPut,
			want:           "updated",
			wantErr:        false,
		},
		{
			name:           "successful DELETE request",
			serverResponse: "deleted",
			serverStatus:   204,
			method:         http.MethodDelete,
			want:           "deleted",
			wantErr:        false,
		},
		{
			name:           "empty response body",
			serverResponse: "",
			serverStatus:   200,
			method:         http.MethodGet,
			want:           "",
			wantErr:        false,
		},
		{
			name:           "large response body",
			serverResponse: strings.Repeat("x", 10000),
			serverStatus:   200,
			method:         http.MethodGet,
			want:           strings.Repeat("x", 10000),
			wantErr:        false,
		},
		{
			name:           "404 not found",
			serverResponse: "not found",
			serverStatus:   404,
			method:         http.MethodGet,
			wantErr:        true,
			errContains:    "404",
		},
		{
			name:           "500 server error",
			serverResponse: "internal error",
			serverStatus:   500,
			method:         http.MethodGet,
			wantErr:        true,
			errContains:    "500",
		},
		{
			name:           "403 forbidden",
			serverResponse: "forbidden",
			serverStatus:   403,
			method:         http.MethodGet,
			wantErr:        true,
			errContains:    "403",
		},
		{
			name:           "299 still success (2xx range)",
			serverResponse: "success",
			serverStatus:   299,
			method:         http.MethodGet,
			want:           "success",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Errorf("Expected method %s, got %s", tt.method, r.Method)
				}
				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Call Download function
			got, err := Download(server.URL, tt.method)

			if (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Error("Download() expected error but got nil")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Download() error = %v, should contain %s", err, tt.errContains)
				}
				return
			}

			if string(got) != tt.want {
				t.Errorf("Download() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestWithCtx(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		serverResponse string
		serverStatus   int
		url            string
		method         string
		want           string
		wantErr        bool
		errContains    string
	}{
		{
			name:           "background context",
			ctx:            context.Background(),
			serverResponse: "response",
			serverStatus:   200,
			method:         http.MethodGet,
			want:           "response",
			wantErr:        false,
		},
		{
			name:           "TODO context",
			ctx:            context.TODO(),
			serverResponse: "todo response",
			serverStatus:   200,
			method:         http.MethodGet,
			want:           "todo response",
			wantErr:        false,
		},
		{
			name:           "cancelled context",
			ctx:            func() context.Context { ctx, cancel := context.WithCancel(context.Background()); cancel(); return ctx }(),
			serverResponse: "should not reach",
			serverStatus:   200,
			method:         http.MethodGet,
			wantErr:        true,
		},
		{
			name:           "timeout context - immediate timeout",
			ctx:            func() context.Context { ctx, cancel := context.WithTimeout(context.Background(), 0); cancel(); return ctx }(),
			serverResponse: "should not reach",
			serverStatus:   200,
			method:         http.MethodGet,
			wantErr:        true,
		},
		{
			name:           "valid context with values",
			ctx:            context.WithValue(context.Background(), "key", "value"),
			serverResponse: "context with values",
			serverStatus:   200,
			method:         http.MethodGet,
			want:           "context with values",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Set URL if not provided
			url := tt.url
			if url == "" {
				url = server.URL
			}

			// Call WithCtx function
			got, err := WithCtx(tt.ctx, url, tt.method)

			if (err != nil) != tt.wantErr {
				t.Errorf("WithCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if string(got) != tt.want {
				t.Errorf("WithCtx() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestWithCtx_InvalidURL(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		url         string
		method      string
		wantErr     bool
		errContains string
	}{
		{
			name:        "empty URL",
			ctx:         context.Background(),
			url:         "",
			method:      http.MethodGet,
			wantErr:     true,
		},
		{
			name:        "invalid URL format",
			ctx:         context.Background(),
			url:         "://invalid",
			method:      http.MethodGet,
			wantErr:     true,
		},
		{
			name:        "missing protocol",
			ctx:         context.Background(),
			url:         "example.com/path",
			method:      http.MethodGet,
			wantErr:     true,
		},
		{
			name:        "unsupported protocol",
			ctx:         context.Background(),
			url:         "ftp://example.com",
			method:      http.MethodGet,
			wantErr:     true,
		},
		{
			name:        "invalid method",
			ctx:         context.Background(),
			url:         "http://example.com",
			method:      "INVALID METHOD\n",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := WithCtx(tt.ctx, tt.url, tt.method)

			if (err != nil) != tt.wantErr {
				t.Errorf("WithCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err == nil {
				t.Error("WithCtx() expected error but got nil")
			}
		})
	}
}

func TestDownload_BinaryData(t *testing.T) {
	// Test that binary data is handled correctly
	binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(binaryData)
	}))
	defer server.Close()

	got, err := Download(server.URL, http.MethodGet)
	if err != nil {
		t.Fatalf("Download() error = %v", err)
	}

	if len(got) != len(binaryData) {
		t.Errorf("Download() length = %d, want %d", len(got), len(binaryData))
	}

	for i := range binaryData {
		if got[i] != binaryData[i] {
			t.Errorf("Download() byte[%d] = %v, want %v", i, got[i], binaryData[i])
		}
	}
}

func TestDownload_UnsupportedMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	// Test with various valid HTTP methods
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodHead,
		http.MethodOptions,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			got, err := Download(server.URL, method)
			if err != nil {
				t.Errorf("Download() with method %s error = %v", method, err)
			}
			if got == nil && method != http.MethodHead {
				t.Error("Download() should return response body")
			}
		})
	}
}

func TestDownload_DownloadErrorWrapper(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	}))
	defer server.Close()

	_, err := Download(server.URL, http.MethodGet)
	if err == nil {
		t.Error("Download() expected error for 400 status")
		return
	}

	// Verify error wraps ErrDownloadFailed
	if !strings.Contains(err.Error(), "Download failed") {
		t.Errorf("Error should contain 'Download failed', got: %v", err)
	}

	if !strings.Contains(err.Error(), "400") {
		t.Errorf("Error should contain status code 400, got: %v", err)
	}
}

func TestDownload_StreamResponse(t *testing.T) {
	// Test with a server that streams data
	chunks := []string{"chunk1", "chunk2", "chunk3"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, _ := w.(http.Flusher)
		for _, chunk := range chunks {
			w.Write([]byte(chunk))
			flusher.Flush()
		}
	}))
	defer server.Close()

	got, err := Download(server.URL, http.MethodGet)
	if err != nil {
		t.Fatalf("Download() error = %v", err)
	}

	expected := strings.Join(chunks, "")
	if string(got) != expected {
		t.Errorf("Download() = %v, want %v", string(got), expected)
	}
}

func TestDownload_ContextCancellation(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if context is cancelled
		select {
		case <-r.Context().Done():
			w.WriteHeader(http.StatusRequestTimeout)
			return
		}
	}))
	defer server.Close()

	// Cancel context immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := WithCtx(ctx, server.URL, http.MethodGet)
	if err == nil {
		t.Error("WithCtx() expected error with cancelled context")
	}
}
