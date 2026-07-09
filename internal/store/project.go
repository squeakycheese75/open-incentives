package store

import (
	"context"

	"github.com/squeakycheese75/open-incentives/internal/db/sqlitedb"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type ProjectStore interface {
	FindByOrgAndPublicID(ctx context.Context, orgID int64, publicID string) (domain.Project, error)
}

type projectStore struct {
	queries *sqlitedb.Queries
}

func (s *projectStore) FindByOrgAndPublicID(ctx context.Context, orgID int64, publicID string) (domain.Project, error) {
	result, err := s.queries.GetProjectByPublicID(ctx, sqlitedb.GetProjectByPublicIDParams{
		PublicID: publicID,
		OrgID:    orgID,
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
