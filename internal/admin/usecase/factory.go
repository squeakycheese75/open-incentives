package usecase_admin

import (
	"context"
	"encoding/json"

	"github.com/squeakycheese75/open-incentives/internal/admin/usecase/campaign_create"
	"github.com/squeakycheese75/open-incentives/internal/admin/usecase/campaign_get"
	"github.com/squeakycheese75/open-incentives/internal/admin/usecase/campaign_list"
	project_create_apikey "github.com/squeakycheese75/open-incentives/internal/admin/usecase/project_apikey_create"
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
	CreateProjectAPIKeyUsecase interface {
		Execute(ctx context.Context, input domain.CreateProjectAPIKEYUsecaseInput) (domain.CreateProjectAPIKEYUsecaseOutput, error)
	}
	ListCampaignsUsecase interface {
		Execute(ctx context.Context, input domain.ListCampaignsUsecaseInput) (domain.ListCampaignsUsecaseOutput, error)
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
	APIKeyStore interface {
		Scope(orgID int64) ScopedAPIKeyStore
	}
	ScopedAPIKeyStore interface {
		Create(ctx context.Context, in domain.APIKey) (domain.APIKey, error)
	}
	CampaignStore interface {
		Scope(orgID int64) ScopedCampaignStore
	}
)

type (
	CryptoSvc interface {
		GenerateKey(size int) ([]byte, error)
	}
	PublicIDGenerator interface {
		New(prefix string) (string, error)
	}
	RuleValidator interface {
		ValidateRules(raw json.RawMessage) error
	}
	PasswordSvc interface {
		Hash(password string) (string, error)
	}
)

type AdminUsecaseFactory struct {
	projectStore  store.ProjectStore
	campaignStore store.CampaignStore
	apiKeyStore   store.APIKeyStore
	idGenerator   PublicIDGenerator
	ruleValidator RuleValidator
	cryptoSvc     CryptoSvc
	passwordSvc   PasswordSvc
}

func NewAdminUsecaseFactory(
	projectStore store.ProjectStore,
	campaignStore store.CampaignStore,
	apiKeyStore store.APIKeyStore,
	idGenerator PublicIDGenerator,
	ruleValidator RuleValidator,
	cryptoSvc CryptoSvc,
	passwordSvc PasswordSvc,
) *AdminUsecaseFactory {
	return &AdminUsecaseFactory{
		projectStore:  projectStore,
		campaignStore: campaignStore,
		apiKeyStore:   apiKeyStore,
		idGenerator:   idGenerator,
		ruleValidator: ruleValidator,
		cryptoSvc:     cryptoSvc,
		passwordSvc:   passwordSvc,
	}
}

func (f *AdminUsecaseFactory) CreateCampaignUsecase(orgID int64) CreateCampaignUsecase {
	return campaign_create.New(f.projectStore.Scope(orgID), f.campaignStore.Scope(orgID), f.idGenerator, f.ruleValidator)
}

func (f *AdminUsecaseFactory) GetCampaignUsecase(orgID int64) GetCampaignUsecase {
	return campaign_get.New(f.campaignStore.Scope(orgID), f.projectStore.Scope(orgID))
}

func (f *AdminUsecaseFactory) ListCampaignUsecase(orgID int64) ListCampaignsUsecase {
	return campaign_list.New(f.campaignStore.Scope(orgID), f.projectStore.Scope(orgID))
}

func (f *AdminUsecaseFactory) CreateProjectAPIKeyUsecase(orgID int64) CreateProjectAPIKeyUsecase {
	return project_create_apikey.New(f.cryptoSvc, f.idGenerator, f.apiKeyStore.Scope(orgID), f.passwordSvc, f.projectStore.Scope(orgID))
}
