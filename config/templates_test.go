package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessTemplates(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "simple template",
			input:    "Hello {{ name }}",
			envVars:  map[string]string{"name": "World"},
			expected: "Hello World",
		},
		{
			name:     "multiple templates",
			input:    "{{ greeting }} {{ name }}!",
			envVars:  map[string]string{"greeting": "Hello", "name": "Go"},
			expected: "Hello Go!",
		},
		{
			name:     "template with spaces",
			input:    "{{ key_with_spaces }}",
			envVars:  map[string]string{"key_with_spaces": "value"},
			expected: "value",
		},
		{
			name:     "template not found",
			input:    "{{ missing }}",
			envVars:  map[string]string{},
			expected: "{{ missing }}",
		},
		{
			name:     "no templates",
			input:    "no templates here",
			envVars:  map[string]string{},
			expected: "no templates here",
		},
		{
			name:     "mixed content",
			input:    "prefix {{ var }} suffix",
			envVars:  map[string]string{"var": "middle"},
			expected: "prefix middle suffix",
		},
		{
			name:     "empty template value",
			input:    "{{ empty }}",
			envVars:  map[string]string{"empty": ""},
			expected: "{{ empty }}",
		},
		{
			name:     "template with underscores",
			input:    "{{ my_auth_token }}",
			envVars:  map[string]string{"my_auth_token": "abc123"},
			expected: "abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			result := ProcessTemplates(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProcessEndpointTemplates(t *testing.T) {
	os.Setenv("AUTH_TOKEN", "bearer-token-123")
	os.Setenv("API_KEY", "key-456")
	defer func() {
		os.Unsetenv("AUTH_TOKEN")
		os.Unsetenv("API_KEY")
	}()

	endpoint := &Endpoint{
		Path:      "/api",
		RemoteURL: "https://api.{{ domain }}.com",
		Headers: []map[string]string{
			{"Authorization": "Bearer {{ AUTH_TOKEN }}"},
			{"X-API-Key": "{{ API_KEY }}"},
		},
		QueryParams: []map[string]string{
			{"token": "{{ AUTH_TOKEN }}"},
			{"version": "v1"},
		},
	}

	os.Setenv("domain", "example")
	defer os.Unsetenv("domain")

	ProcessEndpointTemplates(endpoint)

	assert.Equal(t, "https://api.example.com", endpoint.RemoteURL)
	assert.Equal(t, []map[string]string{
		{"Authorization": "Bearer bearer-token-123"},
		{"X-API-Key": "key-456"},
	}, endpoint.Headers)
	assert.Equal(t, []map[string]string{
		{"token": "bearer-token-123"},
		{"version": "v1"},
	}, endpoint.QueryParams)
}

func TestProcessEndpointTemplatesNoEnvVars(t *testing.T) {
	endpoint := &Endpoint{
		Path:      "/api",
		RemoteURL: "https://api.{{ missing }}.com",
		Headers: []map[string]string{
			{"Authorization": "Bearer {{ missing_token }}"},
		},
		QueryParams: []map[string]string{
			{"key": "{{ missing_key }}"},
		},
	}

	ProcessEndpointTemplates(endpoint)

	assert.Equal(t, "https://api.{{ missing }}.com", endpoint.RemoteURL)
	assert.Equal(t, []map[string]string{
		{"Authorization": "Bearer {{ missing_token }}"},
	}, endpoint.Headers)
	assert.Equal(t, []map[string]string{
		{"key": "{{ missing_key }}"},
	}, endpoint.QueryParams)
}