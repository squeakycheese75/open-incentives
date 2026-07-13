package store

import (
	"context"

	"github.com/squeakycheese75/open-incentives/internal/db/sqlitedb"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type APIKeyStore interface {
	Scope(orgID int64) ScopedAPIKeyStore
	FindByPublicID(ctx context.Context, publicID string) (domain.APIKey, error)
}

type ScopedAPIKeyStore interface {
	Create(ctx context.Context, apikey domain.APIKey) (domain.APIKey, error)
	//
}

type apikeyStore struct {
	queries *sqlitedb.Queries
}

type scopedApikeyStore struct {
	queries *sqlitedb.Queries
	orgID   int64
}

func (s *apikeyStore) Scope(orgID int64) ScopedAPIKeyStore {
	return &scopedApikeyStore{
		queries: s.queries,
		orgID:   orgID,
	}
}

func (s *apikeyStore) FindByPublicID(ctx context.Context, publicID string) (domain.APIKey, error) {
	out, err := s.queries.GetActiveAPIKeyByPublicID(ctx, publicID)
	if err != nil {
		return domain.APIKey{}, err
	}

	return domain.APIKey{
		PublicID:  out.PublicID,
		OrgID:     out.OrgID,
		ProjectID: out.ProjectID,
		KeyHash:   out.KeyHash,
	}, nil
}

func (s *scopedApikeyStore) Create(ctx context.Context, in domain.APIKey) (domain.APIKey, error) {
	out, err := s.queries.CreateProjectAPIKey(ctx, sqlitedb.CreateProjectAPIKeyParams{
		Name:      in.Name,
		PublicID:  in.PublicID,
		OrgID:     s.orgID,
		ProjectID: in.ProjectID,
		KeyHash:   in.KeyHash,
		Prefix:    in.Prefix,
		Status:    string(in.Status),
	})
	if err != nil {
		return domain.APIKey{}, err
	}
	return domain.APIKey{
		Name:      out.Name,
		PublicID:  out.PublicID,
		OrgID:     out.OrgID,
		ProjectID: out.ProjectID,
		KeyHash:   out.KeyHash,
		Prefix:    out.Prefix,
		Status:    domain.APIKeyStatus(out.Status),
		CreatedAt: out.CreatedAt,
	}, nil
}
