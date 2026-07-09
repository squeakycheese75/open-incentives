package httpserver

import (
	"net/http"

	"github.com/squeakycheese75/open-incentives/internal/admin"
	auth "github.com/squeakycheese75/open-incentives/internal/admin/auth"
	middleware "github.com/squeakycheese75/open-incentives/internal/auth"

	"github.com/squeakycheese75/open-incentives/internal/eval"
)

type TokenVerifier interface {
	Verify(tokenString string) (*middleware.Claims, error)
}

func New(
	adminHandler *admin.Handler,
	authHandler *auth.Handler,
	evalHandler *eval.Handler,
	tokenVerifier TokenVerifier,
) http.Handler {
	root := http.NewServeMux()

	adminMux := http.NewServeMux()
	admin.RegisterProtected(adminMux, adminHandler)

	authMux := http.NewServeMux()
	auth.RegisterPublic(authMux, authHandler)

	root.Handle("/admin/auth/",
		http.StripPrefix("/admin/auth", authMux),
	)

	root.Handle("/admin/",
		middleware.AdminAuthMiddleware(tokenVerifier)(
			http.StripPrefix("/admin", adminMux),
		),
	)

	evalMux := http.NewServeMux()
	eval.Register(evalMux, evalHandler)

	root.Handle("/v1/",
		http.StripPrefix("/v1", evalMux),
	)

	return root
}
