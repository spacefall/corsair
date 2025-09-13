package config

import (
	"os"
	"regexp"
	"strings"
)

var templateRegex = regexp.MustCompile(`\{\{\s*([^}]+)\s*\}\}`)

func ProcessTemplates(input string) string {
	return templateRegex.ReplaceAllStringFunc(input, func(match string) string {
		varName := strings.TrimSpace(match[2 : len(match)-2])
		if value := os.Getenv(varName); value != "" {
			return value
		}
		return match
	})
}

func ProcessEndpointTemplates(endpoint *Endpoint) {
	endpoint.RemoteURL = ProcessTemplates(endpoint.RemoteURL)
	
	for i, headerMap := range endpoint.Headers {
		for key, value := range headerMap {
			endpoint.Headers[i][key] = ProcessTemplates(value)
		}
	}
	
	for i, paramMap := range endpoint.QueryParams {
		for key, value := range paramMap {
			endpoint.QueryParams[i][key] = ProcessTemplates(value)
		}
	}
}