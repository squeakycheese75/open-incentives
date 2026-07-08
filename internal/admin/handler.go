package admin

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/squeakycheese75/open-incentives/internal/domain"
)

var ErrInvalidRequest = errors.New("invalid request")

type (
	CampaignStore interface {
		Create(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error)
		Get(ctx context.Context, id string) (domain.Campaign, error)
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

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
