package campaign_update

import (
	"context"
	"fmt"
	"strings"

	"github.com/squeakycheese75/open-incentives/internal/domain"
)

//go:generate mockgen -source=main.go -destination=mocks/mocks.go -package=mocks

type (
	ScopedProjectStore interface {
		Find(ctx context.Context, publicID string) (domain.Project, error)
	}
	ScopedCampaignStore interface {
		Update(ctx context.Context, campaignPublicID string, projectID int64, update domain.Campaign) (domain.Campaign, error)
	}
)

type Usecase struct {
	campaigns ScopedCampaignStore
	projects  ScopedProjectStore
}

func New(campaigns ScopedCampaignStore, projects ScopedProjectStore) *Usecase {
	return &Usecase{campaigns: campaigns, projects: projects}
}

func (uc *Usecase) Execute(ctx context.Context, input domain.UpdateCampaignUsecaseInput) (domain.UpdateCampaignUsecaseOutput, error) {
	campaignPublicID := strings.TrimSpace(input.CampaignPublicID)
	if campaignPublicID == "" {
		return domain.UpdateCampaignUsecaseOutput{}, fmt.Errorf("campaign id is required: %w", domain.ErrInvalidInput)
	}

	projectPublicID := strings.TrimSpace(input.ProjectPublicID)
	if projectPublicID == "" {
		return domain.UpdateCampaignUsecaseOutput{}, fmt.Errorf("project id is required: %w", domain.ErrInvalidInput)
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.UpdateCampaignUsecaseOutput{}, fmt.Errorf("campaign name is required: %w", domain.ErrInvalidInput)
	}

	status := domain.CampaignStatus(strings.TrimSpace(input.Status))
	if status != domain.CampaignStatusActive && status != domain.CampaignStatusInactive {
		return domain.UpdateCampaignUsecaseOutput{}, fmt.Errorf("campaign status must be active or inactive: %w", domain.ErrInvalidInput)
	}

	project, err := uc.projects.Find(ctx, projectPublicID)
	if err != nil {
		return domain.UpdateCampaignUsecaseOutput{}, err
	}

	campaign, err := uc.campaigns.Update(ctx, campaignPublicID, project.ID, domain.Campaign{
		Name:   name,
		Status: status,
		Rules:  input.Rules,
	})
	if err != nil {
		return domain.UpdateCampaignUsecaseOutput{}, err
	}

	return domain.UpdateCampaignUsecaseOutput{
		Campaign: campaign,
	}, nil
}
