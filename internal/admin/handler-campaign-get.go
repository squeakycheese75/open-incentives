package admin

import "net/http"

func (s *Handler) GetCampaign(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	campaign, err := s.store.Get(r.Context(), slug)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "failed_to_create_campaign",
			"msg":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, campaign)
}
