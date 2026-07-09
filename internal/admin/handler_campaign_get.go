package admin

import (
	"net/http"

	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

func (s *Handler) GetCampaign(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusNotImplemented, map[string]any{"error": "not implemented"})
}
