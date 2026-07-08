package httpserver

import (
	"net/http"

	"github.com/squeakycheese75/open-incentives/internal/admin"
	auth "github.com/squeakycheese75/open-incentives/internal/admin/auth"

	"github.com/squeakycheese75/open-incentives/internal/eval"
)

func New(
	adminHandler *admin.Handler,
	authHandler *auth.Handler,
	evalHandler *eval.Handler,
) http.Handler {
	root := http.NewServeMux()

	adminMux := http.NewServeMux()
	admin.RegisterProtected(adminMux, adminHandler)

	adminAuthMux := http.NewServeMux()
	auth.RegisterPublic(adminAuthMux, authHandler)

	evalMux := http.NewServeMux()
	eval.Register(evalMux, evalHandler)

	root.Handle("/admin/",
		http.StripPrefix("/admin", adminMux),
	)

	root.Handle("/v1/",
		http.StripPrefix("/v1", evalMux),
	)

	return root
}

// mux := http.NewServeMux(adminHandler, evalHandler)

// emaxmux.Register(m, "admin", evalHandler)
// adminmux.Register(m, "admin", adminHandler)
// mu

// mux.HandleFunc("GET /admin/health", adminHandler.Health)
// mux.HandleFunc("POST /admin/campaigns", adminHandler.CreateCampaign)
// mux.HandleFunc("GET /admin/campaigns/{slug}", adminHandler.GetCampaign)
// mux.HandleFunc("POST /admin/orgs", adminHandler.CreateOrg)
// POST /orgs
// POST /projects
// POST /projects/{project_id}/api-keys
// POST /projects/{project_id}/campaigns

// mux.HandleFunc("GET /evaluate/health", evalHandler.Health)
// mux.HandleFunc("POST /v1/evaluate", evalHandler.Evaluate)

// GET    /v1/campaigns
// POST   /v1/campaigns
// GET    /v1/campaigns/:id
// PATCH  /v1/campaigns/:id
// DELETE /v1/campaigns/:id
// POST   /v1/evaluate
