package project_update

import (
	"context"
	"fmt"
	"strings"

	"github.com/squeakycheese75/open-incentives/internal/domain"
)

//go:generate mockgen -source=main.go -destination=mocks/mocks.go -package=mocks

type ScopedProjectStore interface {
	Update(ctx context.Context, publicID string, name string) (domain.Project, error)
}

type Usecase struct {
	projects ScopedProjectStore
}

func New(projects ScopedProjectStore) *Usecase {
	return &Usecase{projects: projects}
}

func (uc *Usecase) Execute(ctx context.Context, input domain.UpdateProjectUsecaseInput) (domain.UpdateProjectUsecaseOutput, error) {
	projectPublicID := strings.TrimSpace(input.ProjectPublicID)
	if projectPublicID == "" {
		return domain.UpdateProjectUsecaseOutput{}, fmt.Errorf("project id is required: %w", domain.ErrInvalidInput)
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.UpdateProjectUsecaseOutput{}, fmt.Errorf("project name is required: %w", domain.ErrInvalidInput)
	}

	project, err := uc.projects.Update(ctx, projectPublicID, name)
	if err != nil {
		return domain.UpdateProjectUsecaseOutput{}, err
	}

	return domain.UpdateProjectUsecaseOutput{
		Project: project,
	}, nil
}
