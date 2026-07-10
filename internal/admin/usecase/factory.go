package usecase_admin

import (
	"context"
	"encoding/json"

	"github.com/squeakycheese75/open-incentives/internal/admin/usecase/campaign_create"
	"github.com/squeakycheese75/open-incentives/internal/admin/usecase/campaign_get"
	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/store"
)

type (
	CreateCampaignUsecase interface {
		Execute(ctx context.Context, input domain.CreateCampaignUsecaseInput) (domain.CreateCampaignUsecaseOutput, error)
	}
	GetCampaignUsecase interface {
		Execute(ctx context.Context, input domain.GetCampaignUsecaseInput) (domain.GetCampaignUsecaseOutput, error)
	}
)

type (
	ScopedProjectStore interface {
		Find(ctx context.Context, publicID string) (domain.Campaign, error)
	}
	ProjectStore interface {
		Scope(orgID int64) ScopedCampaignStore
	}
	ScopedCampaignStore interface {
		Create(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error)
		Find(ctx context.Context, publicID string) (domain.Campaign, error)
	}
	CampaignStore interface {
		Scope(orgID int64) ScopedCampaignStore
	}
	PublicIDGenerator interface {
		New(prefix string) (string, error)
	}
	RuleValidator interface {
		ValidateRules(raw json.RawMessage) error
	}
)

type AdminUsecaseFactory struct {
	createCampaignUsecase CreateCampaignUsecase
	getCampaignUsecase    GetCampaignUsecase
}

func NewAdminUsecaseFactory(projectStore store.ProjectStore, campaignStore store.CampaignStore, idGenerator PublicIDGenerator, ruleValidator RuleValidator) *AdminUsecaseFactory {
	return &AdminUsecaseFactory{
		createCampaignUsecase: campaign_create.New(projectStore, campaignStore, idGenerator, ruleValidator),
		getCampaignUsecase:    campaign_get.New(campaignStore, projectStore),
	}
}

func (f *AdminUsecaseFactory) CreateCampaignUsecase() CreateCampaignUsecase {
	return f.createCampaignUsecase
}

func (f *AdminUsecaseFactory) GetCampaignUsecase() GetCampaignUsecase {
	return f.getCampaignUsecase
}
