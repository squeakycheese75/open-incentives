package admin

import (
	"errors"
	"net/http"
	"time"

	"github.com/squeakycheese75/open-incentives/internal/auth"
	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

func (s *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	authCtx, ok := auth.AuthFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "missing auth context")
		return
	}

	out, err := s.adminContainer.ListProjectsUsecase(authCtx.OrgID).Execute(r.Context())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		case errors.Is(err, domain.ErrNotFound):
			httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		default:
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "failed_to_list_projects"})
		}
		return
	}

	projects := make([]ProjectResponse, 0, len(out.Projects))
	for _, p := range out.Projects {
		projects = append(projects, projectToResponse(p))
	}

	httputil.WriteJSON(w, http.StatusOK, ListProjectsResponse{Projects: projects})
}

type ProjectResponse struct {
	PublicID  string `json:"publicId"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ListProjectsResponse struct {
	Projects []ProjectResponse `json:"projects"`
}

func projectToResponse(p domain.Project) ProjectResponse {
	return ProjectResponse{
		PublicID:  p.PublicID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
