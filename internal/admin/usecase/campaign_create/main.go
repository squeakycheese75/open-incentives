package campaign_create

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/store"
)

//	type ProjectStore interface {
//		Find(ctx context.Context, orgID int64, publicID string) (domain.Project, error)
//	}
// type ScopedProjectStore interface {
// 	Find(ctx context.Context, publicID string) (domain.Project, error)
// }

type ProjectStore interface {
	// Scope(orgID int64) interface {
	// 	Find(ctx context.Context, publicID string) (domain.Project, error)
	// }
	Scope(orgID int64) store.ScopedProjectStore
	// Scope(orgID int64) ScopedProjectStore
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
		return domain.CreateCampaignUsecaseOutput{}, errors.New("org id is required")
	}

	projectPublicID := strings.TrimSpace(input.ProjectPublicID)
	if projectPublicID == "" {
		return domain.CreateCampaignUsecaseOutput{}, errors.New("project id is required")
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.CreateCampaignUsecaseOutput{}, errors.New("campaign name is required")
	}

	project, err := uc.projects.Scope(input.OrgID).Find(ctx, projectPublicID)
	if err != nil {
		return domain.CreateCampaignUsecaseOutput{}, err
	}

	ruleJSON, err := json.Marshal(input.Rules)
	if err != nil {
		return domain.CreateCampaignUsecaseOutput{}, errors.New("invalid campaign rules")
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
