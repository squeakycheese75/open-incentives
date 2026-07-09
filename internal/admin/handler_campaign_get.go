package admin

import (
	"net/http"

	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

func (s *Handler) GetCampaign(w http.ResponseWriter, r *http.Request) {
	// slug := r.PathValue("slug")
	// if slug == "" {
	// 	http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
	// 	return
	// }

	// campaign, err := s.campaignStore.Find(r.Context(), slug)
	// if err != nil {
	// 	httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{
	// 		"error": "failed_to_create_campaign",
	// 		"msg":   err.Error(),
	// 	})
	// 	return
	// }

	httputil.WriteJSON(w, http.StatusOK, nil)
}
