package httputil

import (
	"net/http"
	"slices"
)

type CORSConfig struct {
	AllowedOrigins []string
}

func (c CORSConfig) isAllowed(origin string) bool {
	return slices.Contains(c.AllowedOrigins, origin)
}

// CORSMiddleware allows a browser-based frontend (served from a different
// origin) to call this API with an Authorization header. It must wrap
// handlers outside of any auth middleware, since preflight OPTIONS requests
// never carry the Authorization header.
func CORSMiddleware(cfg CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && cfg.isAllowed(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}

			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				w.Header().Set("Access-Control-Max-Age", "600")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
