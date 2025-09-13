package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetEffectiveTimeout(t *testing.T) {
	tests := []struct {
		name            string
		config          Config
		endpoint        Endpoint
		expectedTimeout time.Duration
	}{
		{
			name: "endpoint with specific timeout",
			config: Config{
				Server: ServerConfig{DefaultTimeout: "10s"},
			},
			endpoint: Endpoint{
				Path:    "/test",
				Timeout: "30s",
			},
			expectedTimeout: 30 * time.Second,
		},
		{
			name: "endpoint without timeout uses global default",
			config: Config{
				Server: ServerConfig{DefaultTimeout: "15s"},
			},
			endpoint: Endpoint{
				Path: "/test",
			},
			expectedTimeout: 15 * time.Second,
		},
		{
			name: "endpoint with empty timeout uses global default",
			config: Config{
				Server: ServerConfig{DefaultTimeout: "20s"},
			},
			endpoint: Endpoint{
				Path:    "/test",
				Timeout: "",
			},
			expectedTimeout: 20 * time.Second,
		},
		{
			name: "fallback to 10s on invalid timeout",
			config: Config{
				Server: ServerConfig{DefaultTimeout: "invalid"},
			},
			endpoint: Endpoint{
				Path: "/test",
			},
			expectedTimeout: 10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetEffectiveTimeout(tt.endpoint)
			assert.Equal(t, tt.expectedTimeout, result)
		})
	}
}

func TestValidateTimeoutConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid server and endpoint timeouts",
			config: Config{
				Server: ServerConfig{DefaultTimeout: "10s"},
				Endpoints: []Endpoint{
					{Path: "/test", RemoteURL: "http://example.com", Timeout: "30s"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid various duration formats",
			config: Config{
				Server: ServerConfig{DefaultTimeout: "1m30s"},
				Endpoints: []Endpoint{
					{Path: "/test1", RemoteURL: "http://example.com", Timeout: "500ms"},
					{Path: "/test2", RemoteURL: "http://example.com", Timeout: "2m"},
					{Path: "/test3", RemoteURL: "http://example.com", Timeout: "1h"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid server timeout",
			config: Config{
				Server: ServerConfig{DefaultTimeout: "not-a-duration"},
			},
			wantErr: true,
		},
		{
			name: "invalid endpoint timeout",
			config: Config{
				Server: ServerConfig{DefaultTimeout: "10s"},
				Endpoints: []Endpoint{
					{Path: "/test", RemoteURL: "http://example.com", Timeout: "invalid"},
				},
			},
			wantErr: true,
		},
		{
			name: "empty timeouts are valid",
			config: Config{
				Server: ServerConfig{DefaultTimeout: ""},
				Endpoints: []Endpoint{
					{Path: "/test", RemoteURL: "http://example.com", Timeout: ""},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTimeoutConfig(&tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}