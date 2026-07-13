package admin

import (
	"net/http"
)

func (s *Handler) LisCampaigns(w http.ResponseWriter, r *http.Request) {
	projectSlug := r.PathValue("project_public_id")
	if projectSlug == "" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	// authCtx, ok := auth.AuthFromContext(r.Context())
	// if !ok {
	// 	httputil.WriteError(w, http.StatusUnauthorized, "missing auth context")
	// 	return
	// }
}
