package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/squeakycheese75/open-incentives/internal/db/sqlitedb"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type ProjectStore interface {
	Scope(orgID int64) ScopedProjectStore
}

type ScopedProjectStore interface {
	Create(ctx context.Context, project domain.Project) (domain.Project, error)
	Find(ctx context.Context, publicID string) (domain.Project, error)
	List(ctx context.Context) ([]domain.Project, error)
	Update(ctx context.Context, publicID string, name string) (domain.Project, error)
	Delete(ctx context.Context, publicID string) error
}

type projectStore struct {
	queries *sqlitedb.Queries
}

type scopedProjectStore struct {
	queries *sqlitedb.Queries
	orgID   int64
}

func (s *projectStore) Scope(orgID int64) ScopedProjectStore {
	return &scopedProjectStore{
		queries: s.queries,
		orgID:   orgID,
	}
}

func (s *scopedProjectStore) Create(ctx context.Context, project domain.Project) (domain.Project, error) {
	result, err := s.queries.CreateProject(ctx, sqlitedb.CreateProjectParams{
		PublicID: project.PublicID,
		OrgID:    s.orgID,
		Name:     project.Name,
	})
	if err != nil {
		return domain.Project{}, err
	}

	return domain.Project{
		ID:        result.ID,
		PublicID:  result.PublicID,
		OrgID:     result.OrgID,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (s *scopedProjectStore) Find(ctx context.Context, publicID string) (domain.Project, error) {
	result, err := s.queries.GetProjectByPublicID(ctx, sqlitedb.GetProjectByPublicIDParams{
		PublicID: publicID,
		OrgID:    s.orgID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Project{}, domain.ErrNotFound
		}
		return domain.Project{}, err
	}

	return domain.Project{
		ID:        result.ID,
		PublicID:  result.PublicID,
		OrgID:     result.OrgID,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (s *scopedProjectStore) List(ctx context.Context) ([]domain.Project, error) {
	rows, err := s.queries.ListProjectsByOrg(ctx, s.orgID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Project, 0, len(rows))
	for _, row := range rows {
		result = append(result, domain.Project{
			ID:        row.ID,
			PublicID:  row.PublicID,
			OrgID:     row.OrgID,
			Name:      row.Name,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		})
	}

	return result, nil
}

func (s *scopedProjectStore) Update(ctx context.Context, publicID string, name string) (domain.Project, error) {
	result, err := s.queries.UpdateProject(ctx, sqlitedb.UpdateProjectParams{
		Name:     name,
		PublicID: publicID,
		OrgID:    s.orgID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Project{}, domain.ErrNotFound
		}
		return domain.Project{}, err
	}

	return domain.Project{
		ID:        result.ID,
		PublicID:  result.PublicID,
		OrgID:     result.OrgID,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (s *scopedProjectStore) Delete(ctx context.Context, publicID string) error {
	return s.queries.DeleteProject(ctx, sqlitedb.DeleteProjectParams{
		PublicID: publicID,
		OrgID:    s.orgID,
	})
}
