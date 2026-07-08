package admin

import (
	"net/http"
)

func RegisterProtected(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /campaigns", h.CreateCampaign)
	mux.HandleFunc("GET /campaigns/{slug}", h.GetCampaign)
}
