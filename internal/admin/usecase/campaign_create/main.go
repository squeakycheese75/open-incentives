package campaign_create

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type ProjectStore interface {
	FindByOrgAndPublicID(ctx context.Context, orgID int64, publicID string) (domain.Project, error)
}

type CampaignStore interface {
	Create(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error)
}

type Usecase struct {
	projects  ProjectStore
	campaigns CampaignStore
}

func New(projects ProjectStore, campaigns CampaignStore) *Usecase {
	return &Usecase{
		projects:  projects,
		campaigns: campaigns,
	}
}

func (uc *Usecase) Execute(ctx context.Context, input domain.CreateCampaignUsecaseInput) (domain.CreateCampaignUsecaseOutput, error) {
	if input.OrgID == 0 {
		return domain.CreateCampaignUsecaseOutput{}, fmt.Errorf("org id is required: %w", domain.ErrInvalidInput)
	}

	projectPublicID := strings.TrimSpace(input.ProjectPublicID)
	if projectPublicID == "" {
		return domain.CreateCampaignUsecaseOutput{}, fmt.Errorf("project id is required: %w", domain.ErrInvalidInput)
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.CreateCampaignUsecaseOutput{}, fmt.Errorf("campaign name is required: %w", domain.ErrInvalidInput)
	}

	project, err := uc.projects.FindByOrgAndPublicID(ctx, input.OrgID, projectPublicID)
	if err != nil {
		return domain.CreateCampaignUsecaseOutput{}, err
	}

	ruleJSON, err := json.Marshal(input.Rules)
	if err != nil {
		return domain.CreateCampaignUsecaseOutput{}, fmt.Errorf("invalid campaign rules: %w", domain.ErrInvalidInput)
	}

	slug, err := newCampaignSlug()
	if err != nil {
		return domain.CreateCampaignUsecaseOutput{}, err
	}

	campaign, err := uc.campaigns.Create(ctx, domain.Campaign{
		ProjectID: project.ID,
		OrgId:     input.OrgID,
		Name:      name,
		Slug:      slug,
		Status:    "active",
		Rule:      ruleJSON,
	})
	if err != nil {
		return domain.CreateCampaignUsecaseOutput{}, err
	}

	return domain.CreateCampaignUsecaseOutput{
		CampaignPublicID: campaign.Slug,
	}, nil
}

func newCampaignSlug() (string, error) {
	return gonanoid.Generate("abcdefghijklmnopqrstuvwxyz0123456789", 12)
}
