package admin

import (
	"encoding/json"
	"net/http"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

func (s *Handler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var req CreateCampaignRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid_json",
		})
		return
	}

	if req.Name == "" {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": "missing_name",
		})
		return
	}

	slug, err := NewCampaignSlug()
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "failed_to_generate_slug",
		})
		return
	}

	ruleJSON, err := json.Marshal(req.Rules)
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid_rules",
		})
		return
	}

	campaign, err := s.store.Create(r.Context(), domain.Campaign{
		Name:   req.Name,
		Slug:   slug,
		Status: "active",
		Rule:   ruleJSON,
	})
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "failed_to_create_campaign",
			"msg":   err.Error(),
		})
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, campaign)
}

func NewCampaignSlug() (string, error) {
	return gonanoid.Generate("abcdefghijklmnopqrstuvwxyz0123456789", 12)
}
