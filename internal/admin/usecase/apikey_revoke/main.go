package apikey_revoke

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
	ScopedAPIKeyStore interface {
		Revoke(ctx context.Context, apiKeyPublicID string, projectID int64) (domain.APIKey, error)
	}
)

type Usecase struct {
	apikeys  ScopedAPIKeyStore
	projects ScopedProjectStore
}

func New(apikeys ScopedAPIKeyStore, projects ScopedProjectStore) *Usecase {
	return &Usecase{apikeys: apikeys, projects: projects}
}

func (uc *Usecase) Execute(ctx context.Context, input domain.RevokeAPIKeyUsecaseInput) (domain.RevokeAPIKeyUsecaseOutput, error) {
	apiKeyPublicID := strings.TrimSpace(input.APIKeyPublicID)
	if apiKeyPublicID == "" {
		return domain.RevokeAPIKeyUsecaseOutput{}, fmt.Errorf("api key id is required: %w", domain.ErrInvalidInput)
	}

	projectPublicID := strings.TrimSpace(input.ProjectPublicID)
	if projectPublicID == "" {
		return domain.RevokeAPIKeyUsecaseOutput{}, fmt.Errorf("project id is required: %w", domain.ErrInvalidInput)
	}

	project, err := uc.projects.Find(ctx, projectPublicID)
	if err != nil {
		return domain.RevokeAPIKeyUsecaseOutput{}, err
	}

	apikey, err := uc.apikeys.Revoke(ctx, apiKeyPublicID, project.ID)
	if err != nil {
		return domain.RevokeAPIKeyUsecaseOutput{}, err
	}

	return domain.RevokeAPIKeyUsecaseOutput{
		APIKey: apikey,
	}, nil
}
