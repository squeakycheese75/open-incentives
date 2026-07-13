package store

import (
	"context"
	"database/sql"

	"github.com/squeakycheese75/open-incentives/internal/db/sqlitedb"
)

type Store interface {
	Campaigns() CampaignStore
	Users() UserStore
	Orgs() OrgStore
	Projects() ProjectStore
	APIKeys() APIKeyStore

	WithTx(ctx context.Context, fn func(Store) error) error
	Close() error
}

type store struct {
	db      *sql.DB
	queries *sqlitedb.Queries
}

func New(db *sql.DB) Store {
	return &store{
		db:      db,
		queries: sqlitedb.New(db),
	}
}

func (s *store) Campaigns() CampaignStore {
	return &campaignStore{
		queries: s.queries,
	}
}

func (s *store) Users() UserStore {
	return &userStore{
		queries: s.queries,
	}
}

func (s *store) Orgs() OrgStore {
	return &orgStore{
		queries: s.queries,
	}
}

func (s *store) Projects() ProjectStore {
	return &projectStore{
		queries: s.queries,
	}
}

func (s *store) APIKeys() APIKeyStore {
	return &apikeyStore{
		queries: s.queries,
	}
}

func (s *store) Close() error {
	return s.db.Close()
}
