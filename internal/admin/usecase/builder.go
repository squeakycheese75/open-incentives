package usecase_admin

import (
	"context"

	"github.com/squeakycheese75/open-incentives/internal/admin/usecase/campaign_create"
	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/store"
)

type (
	CreateCampaignUsecase interface {
		Execute(ctx context.Context, input domain.CreateCampaignUsecaseInput) (domain.CreateCampaignUsecaseOutput, error)
	}
)

type ScopedProjectStore interface {
	Find(ctx context.Context, publicID string) (domain.Project, error)
}

type (
	ProjectStore interface {
		// Scope(orgID int64) interface {
		// 	Find(ctx context.Context, publicID string) (domain.Project, error)
		// }
		Scope(orgID int64) store.ScopedProjectStore
		// Scope(orgID int64) ScopedProjectStore
	}
	CampaignStore interface {
		Create(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error)
	}
)

type AdminUsecaseContainer struct {
	projectStore  ProjectStore
	campaignStore CampaignStore
}

func NewAdminUsecaseContainer(projectStore ProjectStore, campaignStore CampaignStore) *AdminUsecaseContainer {
	return &AdminUsecaseContainer{
		projectStore:  projectStore,
		campaignStore: campaignStore,
	}
}

func (uc *AdminUsecaseContainer) CreateContainerUsecase() CreateCampaignUsecase {
	return campaign_create.New(uc.projectStore, uc.campaignStore)
}
