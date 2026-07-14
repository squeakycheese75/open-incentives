package admin

import (
	"net/http"
)

func RegisterProtected(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /projects/{project_public_id}/campaigns", h.CreateCampaign)
	mux.HandleFunc("GET /projects/{project_public_id}/campaigns/{campaign_public_id}", h.GetCampaign)
	mux.HandleFunc("GET /projects/{project_public_id}/campaigns", h.ListCampaigns)
	mux.HandleFunc("DELETE /projects/{project_public_id}/campaigns/{campaign_public_id}", h.DeleteCampaign)
	mux.HandleFunc("PATCH /projects/{project_public_id}/campaigns/{campaign_public_id}", h.UpdateCampaign)

	mux.HandleFunc("POST /projects/{project_public_id}/api-keys", h.CreateProjectAPIKey)
	mux.HandleFunc("GET /projects/{project_public_id}/api-keys", h.ListAPIKeys)
	mux.HandleFunc("POST /projects/{project_public_id}/api-keys/{api_key_public_id}/revoke", h.RevokeProjectAPIKey)

	mux.HandleFunc("POST /projects", h.CreateProject)
	mux.HandleFunc("GET /projects", h.ListProjects)
	mux.HandleFunc("PATCH /projects/{project_public_id}", h.UpdateProject)
	mux.HandleFunc("DELETE /projects/{project_public_id}", h.DeleteProject)

}
