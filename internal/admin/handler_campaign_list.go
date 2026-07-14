package admin

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/squeakycheese75/open-incentives/internal/auth"
	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

func (s *Handler) ListCampaigns(w http.ResponseWriter, r *http.Request) {
	projectSlug := r.PathValue("project_public_id")
	if projectSlug == "" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	authCtx, ok := auth.AuthFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "missing auth context")
		return
	}

	out, err := s.adminContainer.ListCampaignUsecase(authCtx.OrgID).Execute(r.Context(), domain.ListCampaignsUsecaseInput{
		ProjectPublicID: projectSlug,
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

	campaigns := make([]CampaignResponse, 0, len(out.Campaigns))

	for _, campaign := range out.Campaigns {
		campaigns = append(campaigns, CampaignResponse{
			PublicID: campaign.PublicID,
			Name:     campaign.Name,
			Status:   string(campaign.Status),
			Rules:    campaign.Rules,
		})
	}

	httputil.WriteJSON(w, http.StatusOK, ListCampaignResponse{
		Campaigns: campaigns,
	})
}

type ListCampaignResponse struct {
	Campaigns []CampaignResponse `json:"campaigns"`
}

type CampaignResponse struct {
	PublicID string          `json:"publicId"`
	Name     string          `json:"name"`
	Status   string          `json:"status"`
	Rules    json.RawMessage `json:"rules"`
}
