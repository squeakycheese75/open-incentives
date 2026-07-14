package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/squeakycheese75/open-incentives/internal/db/sqlitedb"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type APIKeyStore interface {
	Scope(orgID int64) ScopedAPIKeyStore
	FindByPublicID(ctx context.Context, publicID string) (domain.APIKey, error)
}

type ScopedAPIKeyStore interface {
	Create(ctx context.Context, apikey domain.APIKey) (domain.APIKey, error)
	List(ctx context.Context, projectID int64) ([]domain.APIKey, error)
	Revoke(ctx context.Context, apiKeyPublicID string, projectID int64) (domain.APIKey, error)
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

func (s *scopedApikeyStore) List(ctx context.Context, projectID int64) ([]domain.APIKey, error) {
	rows, err := s.queries.ListAPIKeysByProject(ctx, sqlitedb.ListAPIKeysByProjectParams{
		OrgID:     s.orgID,
		ProjectID: projectID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]domain.APIKey, 0, len(rows))
	for _, row := range rows {
		result = append(result, domain.APIKey{
			Name:       row.Name,
			PublicID:   row.PublicID,
			OrgID:      row.OrgID,
			ProjectID:  row.ProjectID,
			Prefix:     row.Prefix,
			Status:     domain.APIKeyStatus(row.Status),
			CreatedAt:  row.CreatedAt,
			LastUsedAt: nullTimePtr(row.LastUsedAt),
			RevokedAt:  nullTimePtr(row.RevokedAt),
		})
	}

	return result, nil
}

func (s *scopedApikeyStore) Revoke(ctx context.Context, apiKeyPublicID string, projectID int64) (domain.APIKey, error) {
	out, err := s.queries.RevokeAPIKey(ctx, sqlitedb.RevokeAPIKeyParams{
		PublicID:  apiKeyPublicID,
		ProjectID: projectID,
		OrgID:     s.orgID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.APIKey{}, domain.ErrNotFound
		}
		return domain.APIKey{}, err
	}

	return domain.APIKey{
		Name:       out.Name,
		PublicID:   out.PublicID,
		OrgID:      out.OrgID,
		ProjectID:  out.ProjectID,
		Prefix:     out.Prefix,
		Status:     domain.APIKeyStatus(out.Status),
		CreatedAt:  out.CreatedAt,
		LastUsedAt: nullTimePtr(out.LastUsedAt),
		RevokedAt:  nullTimePtr(out.RevokedAt),
	}, nil
}

func nullTimePtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}
