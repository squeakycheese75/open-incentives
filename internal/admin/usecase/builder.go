package usecase_admin

import (
	"context"

	"github.com/squeakycheese75/open-incentives/internal/admin/usecase/campaign_create"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type (
	CreateCampaignUsecase interface {
		Execute(ctx context.Context, input domain.CreateCampaignUsecaseInput) (domain.CreateCampaignUsecaseOutput, error)
	}
)

type (
	ProjectStore interface {
		FindByOrgAndPublicID(ctx context.Context, orgID int64, publicID string) (domain.Project, error)
	}
	CampaignStore interface {
		Create(ctx context.Context, campaign domain.Campaign) (domain.Campaign, error)
	}
)

type AdminUsecaseFactory struct {
	createCampaignUsecase CreateCampaignUsecase
}

func NewAdminUsecaseFactory(projectStore ProjectStore, campaignStore CampaignStore) *AdminUsecaseFactory {
	return &AdminUsecaseFactory{
		createCampaignUsecase: campaign_create.New(projectStore, campaignStore),
	}
}

func (uc *AdminUsecaseFactory) CreateContainerUsecase() CreateCampaignUsecase {
	return uc.createCampaignUsecase
}
