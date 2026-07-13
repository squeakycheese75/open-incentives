package store

import (
	"context"

	"github.com/squeakycheese75/open-incentives/internal/db/sqlitedb"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type ProjectStore interface {
	Scope(orgID int64) ScopedProjectStore
}

type ScopedProjectStore interface {
	Find(ctx context.Context, publicID string) (domain.Project, error)
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

func (s *scopedProjectStore) Find(ctx context.Context, publicID string) (domain.Project, error) {
	result, err := s.queries.GetProjectByPublicID(ctx, sqlitedb.GetProjectByPublicIDParams{
		PublicID: publicID,
		OrgID:    s.orgID,
	})
	if err != nil {
		return domain.Project{}, err
	}

	return domain.Project{
		ID:       result.ID,
		Name:     result.Name,
		PublicID: result.PublicID,
	}, nil
}
