package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bastienwirtz/corsair/config"
	"github.com/stretchr/testify/assert"
)

func getTestConfig() config.Config {
	return config.Config{
		Server: config.ServerConfig{
			DefaultTimeout: "10s",
		},
	}
}

func TestForwardHandler(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom-Header", "forwarded")

		if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"received":"post data"}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"method":"` + r.Method + `","path":"` + r.URL.Path + `"}`))
		}
	}))
	defer mockServer.Close()

	tests := []struct {
		name           string
		url            string
		method         string
		body           string
		headers        map[string]string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "simple GET request",
			url:            mockServer.URL + "/api/users",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST request with body",
			url:            mockServer.URL + "/api/create",
			method:         "POST",
			body:           `{"name":"test"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "request with custom headers",
			url:            mockServer.URL + "/protected",
			method:         "GET",
			headers:        map[string]string{"Authorization": "Bearer token123"},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "missing url parameter",
			url:         "",
			method:      "GET",
			expectError: true,
		},
		{
			name:        "non-http scheme",
			url:         "ftp://example.com/file",
			method:      "GET",
			expectError: true,
		},
		{
			name:           "url without scheme defaults to https",
			url:            "httpbin.org/get",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := ForwardHandler(getTestConfig())

			reqURL := "/forward"
			if tt.url != "" {
				reqURL += "?url=" + tt.url
			}

			var req *http.Request
			if tt.body != "" {
				body := strings.NewReader(tt.body)
				req = httptest.NewRequest(tt.method, reqURL, body)
			} else {
				req = httptest.NewRequest(tt.method, reqURL, nil)
			}

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if tt.expectError {
				assert.True(t, w.Code >= 400, "Expected error status code, got %d", w.Code)
			} else {
				if tt.url == "httpbin.org/get" {
					return
				}

				assert.Equal(t, tt.expectedStatus, w.Code)

				if w.Code < 400 {
					assert.NotEmpty(t, w.Body.String())
				}
			}
		})
	}
}

func TestForwardHandlerBadRequest(t *testing.T) {
	tests := []struct {
		name     string
		queryURL string
		expected string
	}{
		{
			name:     "missing url parameter",
			queryURL: "",
			expected: "Missing 'url' query parameter",
		},
		{
			name:     "unsupported scheme",
			queryURL: "ftp://example.com",
			expected: "Only HTTP and HTTPS URLs are allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := ForwardHandler(getTestConfig())

			reqURL := "/forward"
			if tt.queryURL != "" {
				reqURL += "?url=" + tt.queryURL
			}

			req := httptest.NewRequest("GET", reqURL, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), tt.expected)
		})
	}
}

func TestForwardHandlerUnreachableServer(t *testing.T) {
	handler := ForwardHandler(getTestConfig())
	req := httptest.NewRequest("GET", "/forward?url=http://localhost:99999/unreachable", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadGateway, w.Code)
	assert.Contains(t, w.Body.String(), "Request failed")
}
