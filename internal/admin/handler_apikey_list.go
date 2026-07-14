package admin

import (
	"errors"
	"net/http"
	"time"

	"github.com/squeakycheese75/open-incentives/internal/auth"
	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

func (s *Handler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
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

	out, err := s.adminContainer.ListAPIKeysUsecase(authCtx.OrgID).Execute(r.Context(), domain.ListAPIKeysUsecaseInput{
		ProjectPublicID: projectSlug,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		case errors.Is(err, domain.ErrNotFound):
			httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		default:
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "failed_to_list_api_keys"})
		}
		return
	}

	apikeys := make([]APIKeyResponse, 0, len(out.APIKeys))
	for _, k := range out.APIKeys {
		apikeys = append(apikeys, apiKeyToResponse(k))
	}

	httputil.WriteJSON(w, http.StatusOK, ListAPIKeysResponse{APIKeys: apikeys})
}

type APIKeyResponse struct {
	PublicID   string  `json:"publicId"`
	Name       string  `json:"name"`
	Prefix     string  `json:"prefix"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"createdAt"`
	LastUsedAt *string `json:"lastUsedAt,omitempty"`
	RevokedAt  *string `json:"revokedAt,omitempty"`
}

type ListAPIKeysResponse struct {
	APIKeys []APIKeyResponse `json:"apiKeys"`
}

func apiKeyToResponse(k domain.APIKey) APIKeyResponse {
	resp := APIKeyResponse{
		PublicID:  k.PublicID,
		Name:      k.Name,
		Prefix:    k.Prefix,
		Status:    string(k.Status),
		CreatedAt: k.CreatedAt.UTC().Format(time.RFC3339),
	}

	if k.LastUsedAt != nil {
		formatted := k.LastUsedAt.UTC().Format(time.RFC3339)
		resp.LastUsedAt = &formatted
	}

	if k.RevokedAt != nil {
		formatted := k.RevokedAt.UTC().Format(time.RFC3339)
		resp.RevokedAt = &formatted
	}

	return resp
}
