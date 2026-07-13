package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/squeakycheese75/open-incentives/internal/db/sqlitedb"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type CampaignStore interface {
	Scope(orgID int64) ScopedCampaignStore
}

type ScopedCampaignStore interface {
	Create(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error)
	Find(ctx context.Context, campaignPublicID string, projectID int64) (domain.Campaign, error)
	List(ctx context.Context, projectID int64) ([]domain.Campaign, error)
}

type campaignStore struct {
	queries *sqlitedb.Queries
}

type scopedCampaignStore struct {
	queries *sqlitedb.Queries
	orgID   int64
}

func (s *campaignStore) Scope(orgID int64) ScopedCampaignStore {
	return &scopedCampaignStore{
		queries: s.queries,
		orgID:   orgID,
	}
}

func (s *scopedCampaignStore) Find(ctx context.Context, campaignPublicID string, projectID int64) (domain.Campaign, error) {
	result, err := s.queries.GetCampaign(ctx, sqlitedb.GetCampaignParams{
		PublicID:  campaignPublicID,
		OrgID:     s.orgID,
		ProjectID: projectID,
	})
	if err != nil {
		return domain.Campaign{}, err
	}

	rules := json.RawMessage(result.Rules)
	if !json.Valid(rules) {
		return domain.Campaign{}, fmt.Errorf(
			"campaign %q contains invalid rules JSON",
			result.PublicID,
		)
	}

	return domain.Campaign{
		PublicID: result.PublicID,
		Name:     result.Name,
		Status:   domain.CampaignStatus(result.Status),
		Rules:    rules,
	}, nil
}

func (s *scopedCampaignStore) List(ctx context.Context, projectID int64) ([]domain.Campaign, error) {
	campaigns, err := s.queries.ListActiveCampaignsByProject(ctx, sqlitedb.ListActiveCampaignsByProjectParams{
		ProjectID: projectID,
		OrgID:     s.orgID,
	})
	if err != nil {
		return []domain.Campaign{}, err
	}

	var result []domain.Campaign
	for _, c := range campaigns {
		rules := json.RawMessage(c.Rules)
		if !json.Valid(rules) {
			return []domain.Campaign{}, fmt.Errorf(
				"campaign %q contains invalid rules JSON",
				c.PublicID,
			)
		}

		result = append(result, domain.Campaign{
			PublicID: c.PublicID,
			Name:     c.Name,
			Status:   domain.CampaignStatus(c.Status),
			Rules:    rules,
		})
	}

	return result, nil
}

func (s *scopedCampaignStore) Create(ctx context.Context, c domain.Campaign) (domain.Campaign, error) {
	result, err := s.queries.CreateCampaign(ctx, sqlitedb.CreateCampaignParams{
		PublicID:  c.PublicID,
		ProjectID: c.ProjectID,
		Name:      c.Name,
		Rules:     c.Rules,
		Status:    string(c.Status),
	})
	if err != nil {
		return domain.Campaign{}, err
	}

	return domain.Campaign{
		Name:      result.Name,
		Status:    domain.CampaignStatus(result.Status),
		ProjectID: result.ProjectID,
		PublicID:  result.PublicID,
		Rules:     result.Rules,
	}, err
}
