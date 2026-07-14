package admin

import (
	"errors"
	"net/http"

	"github.com/squeakycheese75/open-incentives/internal/auth"
	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

func (s *Handler) RevokeProjectAPIKey(w http.ResponseWriter, r *http.Request) {
	projectSlug := r.PathValue("project_public_id")
	if projectSlug == "" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	apiKeySlug := r.PathValue("api_key_public_id")
	if apiKeySlug == "" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	authCtx, ok := auth.AuthFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "missing auth context")
		return
	}

	out, err := s.adminContainer.RevokeAPIKeyUsecase(authCtx.OrgID).Execute(r.Context(), domain.RevokeAPIKeyUsecaseInput{
		APIKeyPublicID:  apiKeySlug,
		ProjectPublicID: projectSlug,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		case errors.Is(err, domain.ErrNotFound):
			httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		default:
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "failed_to_revoke_api_key"})
		}
		return
	}

	httputil.WriteJSON(w, http.StatusOK, apiKeyToResponse(out.APIKey))
}
