package store

import (
	"context"

	"github.com/squeakycheese75/open-incentives/internal/db/sqlitedb"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type CampaignStore interface {
	Create(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error)
	Find(ctx context.Context, id string) (domain.Campaign, error)
}

type campaignStore struct {
	queries *sqlitedb.Queries
}

func (s *campaignStore) Find(ctx context.Context, publicID string) (domain.Campaign, error) {
	result, err := s.queries.GetCampaign(ctx, publicID)
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

func (s *campaignStore) Create(ctx context.Context, c domain.Campaign) (domain.Campaign, error) {
	result, err := s.queries.CreateCampaign(ctx, sqlitedb.CreateCampaignParams{
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
