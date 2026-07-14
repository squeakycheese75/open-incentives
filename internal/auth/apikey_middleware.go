package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

type APIKeyStore interface {
	FindByPublicID(ctx context.Context, publicID string) (domain.APIKey, error)
}

type APIKeyVerifier interface {
	Verify(secret, hash string) bool
}

type EvalAuthContext struct {
	APIKeyPublicID string
	ProjectID      int64
	OrgID          int64
}

type evalAuthContextKey struct{}

func ContextWithEvalAuth(ctx context.Context, authCtx EvalAuthContext) context.Context {
	return context.WithValue(ctx, evalAuthContextKey{}, authCtx)
}

func EvalAuthFromContext(ctx context.Context) (EvalAuthContext, bool) {
	authCtx, ok := ctx.Value(evalAuthContextKey{}).(EvalAuthContext)
	return authCtx, ok
}

func EvalAuthMiddleware(
	keyStore APIKeyStore,
	verifier APIKeyVerifier,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				httputil.WriteError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			apiKey, ok := strings.CutPrefix(header, "ApiKey ")
			if !ok || apiKey == "" {
				httputil.WriteError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			publicID, secret, ok := strings.Cut(apiKey, ".")
			if !ok || publicID == "" || secret == "" {
				httputil.WriteError(w, http.StatusUnauthorized, "invalid api key")
				return
			}

			record, err := keyStore.FindByPublicID(r.Context(), publicID)
			if err != nil {
				httputil.WriteError(w, http.StatusUnauthorized, "invalid api key")
				return
			}

			if !verifier.Verify(secret, record.KeyHash) {
				httputil.WriteError(w, http.StatusUnauthorized, "invalid api key")
				return
			}

			authCtx := EvalAuthContext{
				APIKeyPublicID: record.PublicID,
				ProjectID:      record.ProjectID,
				OrgID:          record.OrgID,
			}

			next.ServeHTTP(
				w,
				r.WithContext(ContextWithEvalAuth(r.Context(), authCtx)),
			)
		})
	}
}
