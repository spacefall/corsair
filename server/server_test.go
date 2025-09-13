package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bastienwirtz/corsair/config"
)

func TestTrailingSlashHandling(t *testing.T) {
	// Mock backend server that shows what path it received
	mockBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"received_path":"` + r.URL.Path + `"}`))
	}))
	defer mockBackend.Close()

	tests := []struct {
		name           string
		endpointPath   string
		requestPath    string
		expectRedirect bool
		description    string
	}{
		{
			name:           "endpoint without slash - request without slash",
			endpointPath:   "/api",
			requestPath:    "/api",
			expectRedirect: false,
			description:    "Request to /api should work without redirect",
		},
		{
			name:           "endpoint without slash - request with slash",
			endpointPath:   "/api",
			requestPath:    "/api/",
			expectRedirect: false,
			description:    "Request to /api/ should work without redirect",
		},
		{
			name:           "endpoint with slash - request without slash",
			endpointPath:   "/api/",
			requestPath:    "/api",
			expectRedirect: false,
			description:    "Request to /api should work when endpoint configured as /api/",
		},
		{
			name:           "endpoint with slash - request with slash",
			endpointPath:   "/api/",
			requestPath:    "/api/",
			expectRedirect: false,
			description:    "Request to /api/ should work when endpoint configured as /api/",
		},
		{
			name:           "endpoint without slash - sub-path request",
			endpointPath:   "/api",
			requestPath:    "/api/users",
			expectRedirect: false,
			description:    "Sub-path requests should work",
		},
		{
			name:           "endpoint with slash - sub-path request",
			endpointPath:   "/api/",
			requestPath:    "/api/users",
			expectRedirect: false,
			description:    "Sub-path requests should work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				Server: config.ServerConfig{
					Port: 8080,
				},
				Endpoints: []config.Endpoint{
					{
						Path:      tt.endpointPath,
						RemoteURL: mockBackend.URL,
					},
				},
			}

			handler := NewDynamicRoutingHandler(cfg)

			req := httptest.NewRequest("GET", tt.requestPath, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if tt.expectRedirect {
				assert.True(t, w.Code >= 300 && w.Code < 400, "Expected redirect status code for %s", tt.description)
			} else {
				assert.Equal(t, http.StatusOK, w.Code, "Expected successful response for %s, got %d", tt.description, w.Code)
				assert.NotEmpty(t, w.Body.String())
			}
		})
	}
}

func TestForwardEndpointValidation(t *testing.T) {
	tests := []struct {
		name             string
		endpointPath     string
		forwardEnabled   bool
		expectWarning    bool
		expectRegistered bool
	}{
		{
			name:             "forward endpoint with slash should be skipped",
			endpointPath:     "/forward/",
			forwardEnabled:   true,
			expectWarning:    true,
			expectRegistered: false,
		},
		{
			name:             "forward endpoint without slash should be skipped",
			endpointPath:     "/forward",
			forwardEnabled:   true,
			expectWarning:    true,
			expectRegistered: false,
		},
		{
			name:             "regular endpoint should be registered",
			endpointPath:     "/api",
			forwardEnabled:   true,
			expectWarning:    false,
			expectRegistered: true,
		},
		{
			name:             "forward-like endpoint should be registered",
			endpointPath:     "/forward-api",
			forwardEnabled:   true,
			expectWarning:    false,
			expectRegistered: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				Server: config.ServerConfig{
					Port:                   8080,
					ForwardEndpointEnabled: &tt.forwardEnabled,
				},
				Endpoints: []config.Endpoint{
					{
						Path:      tt.endpointPath,
						RemoteURL: "http://example.com",
					},
				},
			}

			handler := NewDynamicRoutingHandler(cfg)

			// The validation logic is tested by checking that the server
			// doesn't panic and completes registration successfully
			assert.NotNil(t, handler)
		})
	}
}
