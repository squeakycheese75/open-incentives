package campaign_get

import (
	"context"
	"fmt"
	"strings"

	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/store"
)

type Usecase struct {
	campaigns store.CampaignStore
	projects  store.ProjectStore
}

func New(campaigns store.CampaignStore, projects store.ProjectStore) *Usecase {
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

	project, err := uc.projects.Scope(input.OrgID).Find(ctx, projectPublicID)
	if err != nil {
		return domain.GetCampaignUsecaseOutput{}, err
	}

	campaign, err := uc.campaigns.Scope(input.OrgID).Find(ctx, campaignPublicID, project.ID)
	if err != nil {
		return domain.GetCampaignUsecaseOutput{}, err
	}

	return domain.GetCampaignUsecaseOutput{
		Campaign: campaign,
	}, nil
}
