package httpserver

import (
	"net/http"
	"time"

	"github.com/squeakycheese75/open-incentives/internal/admin"
	auth "github.com/squeakycheese75/open-incentives/internal/admin/auth"
	middleware "github.com/squeakycheese75/open-incentives/internal/auth"

	"github.com/squeakycheese75/open-incentives/internal/eval"
)

func New(
	adminHandler *admin.Handler,
	authHandler *auth.Handler,
	evalHandler *eval.Handler,
	tokenVerifier middleware.TokenVerifier,
	adminContextStore middleware.OrgContextStore,
	apiKeyStore middleware.ApiKeyContextStore,
	apiKeyVerifier middleware.APIKeyVerifier,
) http.Handler {
	root := http.NewServeMux()

	adminMux := http.NewServeMux()
	admin.RegisterProtected(adminMux, adminHandler)

	authMux := http.NewServeMux()
	auth.RegisterPublic(authMux, authHandler)

	root.Handle("/admin/auth/",
		http.StripPrefix("/admin/auth", authMux),
	)

	adminChain := middleware.AdminAuthMiddleware(tokenVerifier)(
		middleware.AdminContextMiddleware(
			adminContextStore,
			middleware.NewAuthContextCache(5*time.Minute),
		)(
			http.StripPrefix("/admin", adminMux),
		),
	)

	root.Handle("/admin/", adminChain)

	evalMux := http.NewServeMux()
	eval.Register(evalMux, evalHandler)

	root.Handle("/v1/",
		middleware.EvalAuthMiddleware(apiKeyStore, apiKeyVerifier)(
			http.StripPrefix("/v1", evalMux),
		),
	)

	return root
}
