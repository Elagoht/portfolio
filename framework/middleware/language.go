// Package middleware provides HTTP middleware for the Statigo framework.
package middleware

import (
	gocontext "context"
	"net/http"
)

// LanguageConfig configures the language middleware.
type LanguageConfig struct {
	SkipPaths    []string // Paths to skip (exact match)
	SkipPrefixes []string // Path prefixes to skip
}

// DefaultLanguageConfig returns default configuration.
func DefaultLanguageConfig() LanguageConfig {
	return LanguageConfig{
		SkipPaths:    []string{"/robots.txt", "/sitemap.xml"},
		SkipPrefixes: []string{"/health/", "/static/", "/webhook/", "/api/"},
	}
}

// Language middleware sets English as the language for all requests.
func Language(_ interface{}, config LanguageConfig) func(http.Handler) http.Handler {
	skipPathsMap := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipPathsMap[path] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path

			// Skip for certain paths
			if skipPathsMap[path] {
				next.ServeHTTP(w, r)
				return
			}

			// Skip for certain path prefixes
			for _, prefix := range config.SkipPrefixes {
				if hasPrefix(path, prefix) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Set language in context (always English)
			ctx := setLanguage(r.Context(), "en")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

// setLanguage creates a new context with the language set.
func setLanguage(ctx gocontext.Context, lang string) gocontext.Context {
	return gocontext.WithValue(ctx, contextKey("language"), lang)
}

type contextKey string

const languageKey contextKey = "language"
