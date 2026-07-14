package project_create

import (
	"context"
	"fmt"
	"strings"

	"github.com/squeakycheese75/open-incentives/internal/domain"
)

const prefix string = "proj"

//go:generate mockgen -source=main.go -destination=mocks/mocks.go -package=mocks

type (
	PublicIDGenerator interface {
		New(prefix string) (string, error)
	}
	ScopedProjectStore interface {
		Create(ctx context.Context, project domain.Project) (domain.Project, error)
	}
)

type Usecase struct {
	idGenerator PublicIDGenerator
	projects    ScopedProjectStore
}

func New(idGenerator PublicIDGenerator, projects ScopedProjectStore) *Usecase {
	return &Usecase{idGenerator: idGenerator, projects: projects}
}

func (uc *Usecase) Execute(ctx context.Context, input domain.CreateProjectUsecaseInput) (domain.CreateProjectUsecaseOutput, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.CreateProjectUsecaseOutput{}, fmt.Errorf("project name is required: %w", domain.ErrInvalidInput)
	}

	publicID, err := uc.idGenerator.New(prefix)
	if err != nil {
		return domain.CreateProjectUsecaseOutput{}, err
	}

	project, err := uc.projects.Create(ctx, domain.Project{
		PublicID: publicID,
		Name:     name,
	})
	if err != nil {
		return domain.CreateProjectUsecaseOutput{}, err
	}

	return domain.CreateProjectUsecaseOutput{
		Project: project,
	}, nil
}
