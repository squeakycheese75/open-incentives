package apikey_list

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
		List(ctx context.Context, projectID int64) ([]domain.APIKey, error)
	}
)

type Usecase struct {
	apikeys  ScopedAPIKeyStore
	projects ScopedProjectStore
}

func New(apikeys ScopedAPIKeyStore, projects ScopedProjectStore) *Usecase {
	return &Usecase{apikeys: apikeys, projects: projects}
}

func (uc *Usecase) Execute(ctx context.Context, input domain.ListAPIKeysUsecaseInput) (domain.ListAPIKeysUsecaseOutput, error) {
	projectPublicID := strings.TrimSpace(input.ProjectPublicID)
	if projectPublicID == "" {
		return domain.ListAPIKeysUsecaseOutput{}, fmt.Errorf("project id is required: %w", domain.ErrInvalidInput)
	}

	project, err := uc.projects.Find(ctx, projectPublicID)
	if err != nil {
		return domain.ListAPIKeysUsecaseOutput{}, err
	}

	apikeys, err := uc.apikeys.List(ctx, project.ID)
	if err != nil {
		return domain.ListAPIKeysUsecaseOutput{}, err
	}

	return domain.ListAPIKeysUsecaseOutput{
		APIKeys: apikeys,
	}, nil
}
