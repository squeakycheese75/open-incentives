package project_delete

import (
	"context"
	"fmt"
	"strings"

	"github.com/squeakycheese75/open-incentives/internal/domain"
)

//go:generate mockgen -source=main.go -destination=mocks/mocks.go -package=mocks

type ScopedProjectStore interface {
	Delete(ctx context.Context, publicID string) error
}

type Usecase struct {
	projects ScopedProjectStore
}

func New(projects ScopedProjectStore) *Usecase {
	return &Usecase{projects: projects}
}

func (uc *Usecase) Execute(ctx context.Context, input domain.DeleteProjectUsecaseInput) error {
	projectPublicID := strings.TrimSpace(input.ProjectPublicID)
	if projectPublicID == "" {
		return fmt.Errorf("project id is required: %w", domain.ErrInvalidInput)
	}

	return uc.projects.Delete(ctx, projectPublicID)
}
