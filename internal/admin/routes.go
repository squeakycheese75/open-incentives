package admin

import (
	"net/http"
)

func RegisterProtected(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /projects/{project_slug}/campaigns", h.CreateCampaign)
	mux.HandleFunc("GET /projects/{project_slug}/campaigns/{campaign_slug}", h.GetCampaign)
}
