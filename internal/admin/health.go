package admin

import (
	"net/http"

	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

func (s *Handler) Health(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
