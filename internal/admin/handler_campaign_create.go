package admin

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/squeakycheese75/open-incentives/internal/auth"
	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

func (s *Handler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	projectSlug := r.PathValue("project_public_id")
	if projectSlug == "" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

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

	authCtx, ok := auth.AuthFromContext(r.Context())
	if !ok {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid_claims",
		})
		return
	}

	out, err := s.adminContainer.CreateCampaignUsecase(authCtx.OrgID).Execute(r.Context(), domain.CreateCampaignUsecaseInput{
		ProjectPublicID: projectSlug,
		Name:            req.Name,
		Rules:           req.Rules,
	})
	if err != nil {
		slog.ErrorContext(
			r.Context(),
			"create campaign failed",
			"error", err,
		)

		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		case errors.Is(err, domain.ErrNotFound):
			httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		default:
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "failed_to_create_campaign"})
		}
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, CreateCampaignResponse{
		CampaignPublicID: out.CampaignPublicID,
	})
}

type CreateCampaignResponse struct {
	CampaignPublicID string `json:"campaign_public_id"`
}
