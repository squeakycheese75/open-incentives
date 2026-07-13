package store

import (
	"context"

	"github.com/squeakycheese75/open-incentives/internal/db/sqlitedb"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type UserStore interface {
	Find(ctx context.Context, email string, orgID int64) (domain.User, error)
}

type userStore struct {
	queries *sqlitedb.Queries
}

func (s *userStore) Find(ctx context.Context, email string, orgID int64) (domain.User, error) {
	result, err := s.queries.GetUserByEmailAndOrg(ctx, sqlitedb.GetUserByEmailAndOrgParams{
		Email: email,
		OrgID: orgID,
	})
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:             result.ID,
		PublicID:       result.PublicID,
		Email:          result.Email,
		HashedPassword: result.PasswordHash,
	}, err
}
