package auth

import (
	"net/http"
)

func RegisterPublic(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /login", h.AuthLogin)
	// mux.HandleFunc("POST /logout", h.UserLogout)
	// mux.HandleFunc("GET /logout", h.Me)
}
