package store

import (
	"context"

	"github.com/squeakycheese75/open-incentives/internal/db/sqlitedb"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type OrgStore interface {
	Find(ctx context.Context, publicID string) (domain.Organization, error)
}

type orgStore struct {
	queries *sqlitedb.Queries
}

func (s *orgStore) Find(ctx context.Context, publicID string) (domain.Organization, error) {
	result, err := s.queries.GetOrgByPublicID(ctx, publicID)
	if err != nil {
		return domain.Organization{}, err
	}

	return domain.Organization{
		ID:       result.ID,
		Name:     result.Name,
		PublicID: result.PublicID,
	}, err
}
