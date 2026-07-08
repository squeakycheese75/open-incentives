package eval

import (
	"encoding/json"
	"net/http"

	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

func (s *Handler) Evaluate(w http.ResponseWriter, r *http.Request) {
	var httpReq EvaluateRequest

	if err := json.NewDecoder(r.Body).Decode(&httpReq); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid_json",
		})
		return
	}

	// need to load rules from the campaign
	engineReq := toEngineRequest(httpReq)

	result, err := s.engine.Evaluate(r.Context(), engineReq)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "evaluation_failed",
		})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, toEvaluateResponse(result))
}
