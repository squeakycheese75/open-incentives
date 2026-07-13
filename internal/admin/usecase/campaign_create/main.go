package campaign_create

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type (
	PublicIDGenerator interface {
		New(prefix string) (string, error)
	}
	CampaignStore interface {
		Create(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error)
	}
	ProjectStore interface {
		Find(ctx context.Context, publicID string) (domain.Project, error)
	}
	RuleValidator interface {
		ValidateRules(raw json.RawMessage) error
	}
	ScopedProjectStore interface {
		Find(ctx context.Context, publicID string) (domain.Project, error)
	}
	ScopedCampaignStore interface {
		Create(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error)
		Find(ctx context.Context, campaignPublicID string, projectID int64) (domain.Campaign, error)
	}
)

type Usecase struct {
	projects      ScopedProjectStore
	campaigns     ScopedCampaignStore
	idGenerator   PublicIDGenerator
	ruleValidator RuleValidator
}

func New(projects ScopedProjectStore, campaigns ScopedCampaignStore, idGenerator PublicIDGenerator, ruleValidator RuleValidator) *Usecase {
	return &Usecase{
		projects:      projects,
		campaigns:     campaigns,
		idGenerator:   idGenerator,
		ruleValidator: ruleValidator,
	}
}

func (uc *Usecase) Execute(ctx context.Context, input domain.CreateCampaignUsecaseInput) (domain.CreateCampaignUsecaseOutput, error) {
	projectPublicID := strings.TrimSpace(input.ProjectPublicID)
	if projectPublicID == "" {
		return domain.CreateCampaignUsecaseOutput{}, fmt.Errorf("project id is required: %w", domain.ErrInvalidInput)
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.CreateCampaignUsecaseOutput{}, fmt.Errorf("campaign name is required: %w", domain.ErrInvalidInput)
	}

	project, err := uc.projects.Find(ctx, projectPublicID)
	if err != nil {
		return domain.CreateCampaignUsecaseOutput{}, err
	}

	campaignPublicID, err := uc.idGenerator.New("camp")
	if err != nil {
		return domain.CreateCampaignUsecaseOutput{}, err
	}

	if err := uc.ruleValidator.ValidateRules(input.Rules); err != nil {
		return domain.CreateCampaignUsecaseOutput{},
			fmt.Errorf("invalid campaign rules: %w: %v", domain.ErrInvalidInput, err)
	}

	campaign, err := uc.campaigns.Create(ctx, domain.Campaign{
		ProjectID: project.ID,
		Name:      name,
		PublicID:  campaignPublicID,
		Status:    domain.CampaignStatusActive,
		Rules:     input.Rules,
	})
	if err != nil {
		return domain.CreateCampaignUsecaseOutput{}, err
	}

	return domain.CreateCampaignUsecaseOutput{
		CampaignPublicID: campaign.PublicID,
	}, nil
}
