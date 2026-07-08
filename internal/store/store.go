package store

import (
	"context"

	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type CampaignStore interface {
	Create(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error)
	Get(ctx context.Context, id string) (domain.Campaign, error)
	WithTx(ctx context.Context, fn func(CampaignStore) error) error
	Close() error
}
