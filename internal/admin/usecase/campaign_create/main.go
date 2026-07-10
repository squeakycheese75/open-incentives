package campaign_create

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/store"
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
)

type Usecase struct {
	projects      store.ProjectStore
	campaigns     store.CampaignStore
	idGenerator   PublicIDGenerator
	ruleValidator RuleValidator
}

func New(projects store.ProjectStore, campaigns store.CampaignStore, idGenerator PublicIDGenerator, ruleValidator RuleValidator) *Usecase {
	return &Usecase{
		projects:      projects,
		campaigns:     campaigns,
		idGenerator:   idGenerator,
		ruleValidator: ruleValidator,
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

	project, err := uc.projects.Scope(input.OrgID).Find(ctx, projectPublicID)
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

	scopedCampaigns := uc.campaigns.Scope(input.OrgID)
	campaign, err := scopedCampaigns.Create(ctx, domain.Campaign{
		ProjectID: project.ID,
		OrgId:     input.OrgID,
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
