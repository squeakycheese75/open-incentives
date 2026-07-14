package admin

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/squeakycheese75/open-incentives/internal/auth"
	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

func (s *Handler) GetCampaign(w http.ResponseWriter, r *http.Request) {
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

	authCtx, ok := auth.AuthFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "missing auth context")
		return
	}

	out, err := s.adminContainer.GetCampaignUsecase(authCtx.OrgID).Execute(r.Context(), domain.GetCampaignUsecaseInput{
		CampaignPublicID: campaignSlug,
		ProjectPublicID:  projectSlug,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		case errors.Is(err, domain.ErrNotFound):
			httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		default:
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "failed_to_get_campaign"})
		}
		return
	}

	campaign := out.Campaign

	httputil.WriteJSON(w, http.StatusOK, GetCampaignResponse{
		PublicID: campaign.PublicID,
		Name:     campaign.Name,
		Status:   string(campaign.Status),
		Rules:    json.RawMessage(campaign.Rules),
	})
}

type GetCampaignResponse struct {
	PublicID string          `json:"publicId"`
	Name     string          `json:"name"`
	Status   string          `json:"status"`
	Rules    json.RawMessage `json:"rules"`
}
