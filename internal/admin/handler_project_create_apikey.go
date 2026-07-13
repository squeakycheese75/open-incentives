package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/squeakycheese75/open-incentives/internal/auth"
	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

type CreateProjectAPIKeyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateProjectAPIKeyResponse struct {
	APIKey         string `json:"apikey"`
	APIKeyPublicID string `json:"apikey_public_id"`
	CreatedAt      string `json:"created_at"`
}

func (s *Handler) CreateProjectAPIKey(w http.ResponseWriter, r *http.Request) {
	projectSlug := r.PathValue("project_public_id")
	if projectSlug == "" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	var req CreateProjectAPIKeyRequest
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

	out, err := s.adminContainer.CreateProjectAPIKeyUsecase().Execute(r.Context(), domain.CreateProjectAPIKEYUsecaseInput{
		OrgID:           authCtx.OrgID,
		ProjectPublicID: projectSlug,
		Name:            req.Name,
		Description:     req.Description,
	})
	if err != nil {
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

	httputil.WriteJSON(w, http.StatusCreated, CreateProjectAPIKeyResponse{
		APIKey:         out.APIKey,
		APIKeyPublicID: out.APIKeyPublicID,
		CreatedAt:      out.CreatedAt.UTC().Format(time.RFC3339),
	})
}
