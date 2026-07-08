package store

import (
	"context"
	"database/sql"

	"github.com/squeakycheese75/open-incentives/internal/db/sqlitedb"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type StoreSqlite struct {
	*sqlitedb.Queries
	db *sql.DB
}

func New(db *sql.DB) *StoreSqlite {
	return &StoreSqlite{
		db:      db,
		Queries: sqlitedb.New(db),
	}
}

func (s *StoreSqlite) Get(ctx context.Context, publicID string) (domain.Campaign, error) {
	result, err := s.Queries.GetCampaign(ctx, publicID)
	if err != nil {
		return domain.Campaign{}, err
	}

	return domain.Campaign{
		Name:   result.Name,
		Status: domain.CampaignStatus(result.Status),
		Slug:   publicID,
		Rule:   result.Rule,
	}, err
}

func (s *StoreSqlite) Create(ctx context.Context, c domain.Campaign) (domain.Campaign, error) {
	result, err := s.Queries.CreateCampaign(ctx, sqlitedb.CreateCampaignParams{
		PublicID: c.Slug,
		Name:     c.Name,
		Rule:     c.Rule,
		Status:   string(c.Status),
	})
	if err != nil {
		return domain.Campaign{}, err
	}

	return domain.Campaign{
		Name:   result.Name,
		Status: domain.CampaignStatus(result.Status),
		Slug:   result.PublicID,
		Rule:   result.Rule,
	}, err
}

func (s *StoreSqlite) WithTx(ctx context.Context, fn func(CampaignStore) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txStore := &StoreSqlite{
		db:      s.db,
		Queries: s.Queries.WithTx(tx),
	}

	if err := fn(txStore); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *StoreSqlite) Close() error {
	return s.db.Close()
}
