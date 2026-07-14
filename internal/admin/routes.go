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

	mux.HandleFunc("POST /projects/{project_public_id}/api-keys", h.CreateProjectAPIKey)

}
