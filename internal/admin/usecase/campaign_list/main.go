package campaign_list

import (
	"context"
	"fmt"
	"strings"

	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type (
	ScopedProjectStore interface {
		Find(ctx context.Context, publicID string) (domain.Project, error)
	}
	ScopedCampaignStore interface {
		List(ctx context.Context, projectID int64) ([]domain.Campaign, error)
	}
)

type Usecase struct {
	campaigns ScopedCampaignStore
	projects  ScopedProjectStore
}

func New(campaigns ScopedCampaignStore, projects ScopedProjectStore) *Usecase {
	return &Usecase{campaigns: campaigns, projects: projects}
}

func (uc *Usecase) Execute(ctx context.Context, input domain.ListCampaignsUsecaseInput) (domain.ListCampaignsUsecaseOutput, error) {
	projectPublicID := strings.TrimSpace(input.ProjectPublicID)
	if projectPublicID == "" {
		return domain.ListCampaignsUsecaseOutput{}, fmt.Errorf("project id is required: %w", domain.ErrInvalidInput)
	}

	project, err := uc.projects.Find(ctx, projectPublicID)
	if err != nil {
		return domain.ListCampaignsUsecaseOutput{}, err
	}

	out, err := uc.campaigns.List(ctx, project.ID)
	if err != nil {
		return domain.ListCampaignsUsecaseOutput{}, err
	}

	var results []domain.Campaign

	for _, c := range out {
		results = append(results, domain.Campaign{
			PublicID: c.PublicID,
			Name:     c.Name,
			KeyHash:  c.KeyHash,
			Status:   domain.CampaignStatus(c.Status),
			Rules:    c.Rules,
		})
	}

	return domain.ListCampaignsUsecaseOutput{
		Campaigns: results,
	}, nil
}
