package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrailingSlash(t *testing.T) {
	// Create a test handler that captures the request path
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Path", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	middleware := TrailingSlash()
	handler := middleware(testHandler)

	tests := []struct {
		name         string
		requestPath  string
		expectedPath string
	}{
		{
			name:         "path without trailing slash gets normalized",
			requestPath:  "/api",
			expectedPath: "/api/",
		},
		{
			name:         "path with trailing slash remains unchanged",
			requestPath:  "/api/",
			expectedPath: "/api/",
		},
		{
			name:         "root path remains unchanged",
			requestPath:  "/",
			expectedPath: "/",
		},
		{
			name:         "sub-path without trailing slash gets normalized",
			requestPath:  "/api/v1/users",
			expectedPath: "/api/v1/users/",
		},
		{
			name:         "sub-path with trailing slash remains unchanged",
			requestPath:  "/api/v1/users/",
			expectedPath: "/api/v1/users/",
		},
		{
			name:         "file path remains unchanged",
			requestPath:  "/static/style.css",
			expectedPath: "/static/style.css",
		},
		{
			name:         "file path in subdirectory remains unchanged",
			requestPath:  "/assets/images/logo.png",
			expectedPath: "/assets/images/logo.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.requestPath, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, tt.expectedPath, w.Header().Get("X-Request-Path"))
		})
	}
}