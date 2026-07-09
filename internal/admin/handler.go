package admin

import (
	"context"
	"errors"

	"github.com/squeakycheese75/open-incentives/internal/domain"
)

var ErrInvalidRequest = errors.New("invalid request")

type (
	CampaignStore interface {
		Create(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error)
		Find(ctx context.Context, id string) (domain.Campaign, error)
	}
)

type Handler struct {
	store CampaignStore
}

func NewHandler(store CampaignStore) *Handler {
	return &Handler{
		store: store,
	}
}
