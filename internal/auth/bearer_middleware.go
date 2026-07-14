package auth

import (
	"net/http"
	"strings"

	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

type TokenVerifier interface {
	Verify(tokenString string) (*Claims, error)
}

func AdminAuthMiddleware(verifier TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				httputil.WriteError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || token == "" {
				httputil.WriteError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			claims, err := verifier.Verify(token)
			if err != nil {
				httputil.WriteError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := ContextWithClaims(r.Context(), *claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
