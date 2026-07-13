package eval

import (
	"net/http"
)

func Register(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /evaluate", h.Evaluate)
}
