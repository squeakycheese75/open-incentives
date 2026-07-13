package campaign_get

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
		Find(ctx context.Context, campaignPublicID string, projectID int64) (domain.Campaign, error)
	}
)

type Usecase struct {
	campaigns ScopedCampaignStore
	projects  ScopedProjectStore
}

func New(campaigns ScopedCampaignStore, projects ScopedProjectStore) *Usecase {
	return &Usecase{campaigns: campaigns, projects: projects}
}

func (uc *Usecase) Execute(ctx context.Context, input domain.GetCampaignUsecaseInput) (domain.GetCampaignUsecaseOutput, error) {
	campaignPublicID := strings.TrimSpace(input.CampaignPublicID)
	if campaignPublicID == "" {
		return domain.GetCampaignUsecaseOutput{}, fmt.Errorf("campaign id is required: %w", domain.ErrInvalidInput)
	}

	projectPublicID := strings.TrimSpace(input.ProjectPublicID)
	if projectPublicID == "" {
		return domain.GetCampaignUsecaseOutput{}, fmt.Errorf("project id is required: %w", domain.ErrInvalidInput)
	}

	project, err := uc.projects.Find(ctx, projectPublicID)
	if err != nil {
		return domain.GetCampaignUsecaseOutput{}, err
	}

	campaign, err := uc.campaigns.Find(ctx, campaignPublicID, project.ID)
	if err != nil {
		return domain.GetCampaignUsecaseOutput{}, err
	}

	return domain.GetCampaignUsecaseOutput{
		Campaign: campaign,
	}, nil
}
