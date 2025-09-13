package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bastienwirtz/corsair/config"
)

func TestCORS(t *testing.T) {
	tests := []struct {
		name           string
		corsConfig     config.CORSConfig
		requestOrigin  string
		requestMethod  string
		expectedOrigin string
		expectedStatus int
	}{
		{
			name: "wildcard origin without credentials",
			corsConfig: config.CORSConfig{
				Origins:     []string{"*"},
				Methods:     "GET, POST, OPTIONS",
				Headers:     "*",
				Credentials: false,
			},
			requestOrigin:  "http://example.com",
			requestMethod:  "GET",
			expectedOrigin: "*",
			expectedStatus: http.StatusOK,
		},
		{
			name: "wildcard origin with credentials blocked",
			corsConfig: config.CORSConfig{
				Origins:     []string{"*"},
				Methods:     "GET, POST, OPTIONS",
				Headers:     "*",
				Credentials: true,
			},
			requestOrigin:  "http://example.com",
			requestMethod:  "GET",
			expectedOrigin: "",
			expectedStatus: http.StatusOK,
		},
		{
			name: "specific origin allowed",
			corsConfig: config.CORSConfig{
				Origins:     []string{"http://localhost:3000", "https://example.com"},
				Methods:     "GET, POST",
				Headers:     "Content-Type, Authorization",
				Credentials: false,
			},
			requestOrigin:  "http://localhost:3000",
			requestMethod:  "GET",
			expectedOrigin: "http://localhost:3000",
			expectedStatus: http.StatusOK,
		},
		{
			name: "origin not allowed",
			corsConfig: config.CORSConfig{
				Origins:     []string{"https://allowed.com"},
				Methods:     "GET, POST",
				Headers:     "*",
				Credentials: false,
			},
			requestOrigin:  "http://malicious.com",
			requestMethod:  "GET",
			expectedOrigin: "",
			expectedStatus: http.StatusOK,
		},
		{
			name: "preflight OPTIONS request",
			corsConfig: config.CORSConfig{
				Origins:     []string{"http://localhost:3000"},
				Methods:     "GET, POST, PUT",
				Headers:     "Content-Type",
				Credentials: true,
			},
			requestOrigin:  "http://localhost:3000",
			requestMethod:  "OPTIONS",
			expectedOrigin: "http://localhost:3000",
			expectedStatus: http.StatusOK,
		},
		{
			name: "subdomain wildcard",
			corsConfig: config.CORSConfig{
				Origins: []string{"*.example.com"},
				Methods: "GET",
				Headers: "*",
			},
			requestOrigin:  "https://api.example.com",
			requestMethod:  "GET",
			expectedOrigin: "https://api.example.com",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			corsHandler := CORS(tt.corsConfig)(handler)

			req := httptest.NewRequest(tt.requestMethod, "/test", nil)
			if tt.requestOrigin != "" {
				req.Header.Set("Origin", tt.requestOrigin)
			}

			w := httptest.NewRecorder()
			corsHandler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedOrigin != "" {
				assert.Equal(t, tt.expectedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
			} else {
				origin := w.Header().Get("Access-Control-Allow-Origin")
				assert.True(t, origin == "" || origin == tt.requestOrigin,
					"Expected empty or matching origin, got: %s", origin)
			}

			assert.Equal(t, tt.corsConfig.Methods, w.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, tt.corsConfig.Headers, w.Header().Get("Access-Control-Allow-Headers"))

			if tt.corsConfig.Credentials {
				assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
			} else {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
			}

			if tt.requestMethod == "OPTIONS" {
				assert.Equal(t, "86400", w.Header().Get("Access-Control-Max-Age"))
			}
		})
	}
}

func TestIsOriginAllowed(t *testing.T) {
	tests := []struct {
		name           string
		origin         string
		allowedOrigins []string
		expected       bool
	}{
		{
			name:           "exact match",
			origin:         "https://example.com",
			allowedOrigins: []string{"https://example.com", "https://other.com"},
			expected:       true,
		},
		{
			name:           "wildcard",
			origin:         "https://example.com",
			allowedOrigins: []string{"*"},
			expected:       true,
		},
		{
			name:           "subdomain wildcard match",
			origin:         "https://api.example.com",
			allowedOrigins: []string{"*.example.com"},
			expected:       true,
		},
		{
			name:           "subdomain wildcard root match",
			origin:         "https://example.com",
			allowedOrigins: []string{"*.example.com"},
			expected:       false,
		},
		{
			name:           "no match",
			origin:         "https://malicious.com",
			allowedOrigins: []string{"https://example.com", "*.trusted.com"},
			expected:       false,
		},
		{
			name:           "empty origin",
			origin:         "",
			allowedOrigins: []string{"*"},
			expected:       false,
		},
		{
			name:           "subdomain no match",
			origin:         "https://api.other.com",
			allowedOrigins: []string{"*.example.com"},
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOriginAllowed(tt.origin, tt.allowedOrigins)
			assert.Equal(t, tt.expected, result)
		})
	}
}
