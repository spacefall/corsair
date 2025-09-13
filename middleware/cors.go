package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/bastienwirtz/corsair/config"
)

// CORS returns a middleware that handles CORS headers and preflight requests.
// Supports wildcard origins (*), specific origins, and subdomain wildcards (*.domain.com).
// Prevents wildcard origins when credentials are enabled for security compliance.
func CORS(corsConfig config.CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			slog.Debug("Processing CORS request",
				"origin", origin,
				"method", r.Method,
				"path", r.URL.Path,
				"configured_origins", corsConfig.Origins,
				"credentials", corsConfig.Credentials)

			// Determine allowed origin: "*" only when single wildcard without credentials
			allowedOrigin := ""
			if corsConfig.WildcardOriginAllowed() && !corsConfig.Credentials {
				allowedOrigin = "*"
			} else if isOriginAllowed(origin, corsConfig.Origins) {
				allowedOrigin = origin
			} else if origin != "" {
				slog.Warn("Origin rejected",
					"origin", origin,
					"allowed_origins", corsConfig.Origins)
			}

			if allowedOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			}

			if corsConfig.Methods != "" {
				w.Header().Set("Access-Control-Allow-Methods", corsConfig.Methods)
			}

			if corsConfig.Headers != "" {
				w.Header().Set("Access-Control-Allow-Headers", corsConfig.Headers)
			}

			if corsConfig.Credentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Handle preflight OPTIONS requests
			if r.Method == "OPTIONS" {
				slog.Debug("Skiping remote url request for preflight request",
					"origin", origin,
					"allowed", allowedOrigin != "",
					"path", r.URL.Path)
				w.Header().Set("Access-Control-Max-Age", "86400")
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// isOriginAllowed checks if origin matches any allowed pattern.
// Supports exact matches, "*" wildcard, and "*.domain.com" subdomain wildcards.
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}

	for _, allowed := range allowedOrigins {
		// Exact match or global wildcard
		if allowed == "*" || allowed == origin {
			slog.Debug("Origin matched", "origin", origin, "pattern", allowed, "match_type", "exact")
			return true
		}

		// Subdomain wildcard: *.domain.com matches sub.domain.com but not domain.com
		if strings.HasPrefix(allowed, "*.") {
			domain := allowed[2:]
			if strings.HasSuffix(origin, "."+domain) {
				slog.Debug("Origin matched", "origin", origin, "pattern", allowed, "match_type", "subdomain_wildcard")
				return true
			}
		}
	}
	slog.Debug("Origin not matched", "origin", origin, "allowed_patterns", allowedOrigins)
	return false
}
