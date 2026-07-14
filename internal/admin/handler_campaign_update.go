package admin

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/squeakycheese75/open-incentives/internal/auth"
	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

func (s *Handler) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	projectSlug := r.PathValue("project_public_id")
	if projectSlug == "" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	campaignSlug := r.PathValue("campaign_public_id")
	if campaignSlug == "" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	var req UpdateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid_json",
		})
		return
	}

	authCtx, ok := auth.AuthFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "missing auth context")
		return
	}

	out, err := s.adminContainer.UpdateCampaignUsecase(authCtx.OrgID).Execute(r.Context(), domain.UpdateCampaignUsecaseInput{
		CampaignPublicID: campaignSlug,
		ProjectPublicID:  projectSlug,
		Name:             req.Name,
		Status:           req.Status,
		Rules:            req.Rules,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		case errors.Is(err, domain.ErrNotFound):
			httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		default:
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "failed_to_update_campaign"})
		}
		return
	}

	httputil.WriteJSON(w, http.StatusOK, UpdateCampaignResponse{
		PublicID: out.PublicID,
		Name:     out.Name,
		Status:   string(out.Status),
		Rules:    out.Rules,
	})
}

type UpdateCampaignRequest struct {
	Name   string          `json:"name"`
	Status string          `json:"status"`
	Rules  json.RawMessage `json:"rules"`
}

type UpdateCampaignResponse struct {
	PublicID string          `json:"publicId"`
	Name     string          `json:"name"`
	Status   string          `json:"status"`
	Rules    json.RawMessage `json:"rules"`
}
