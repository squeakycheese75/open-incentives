package project_list

import (
	"context"

	"github.com/squeakycheese75/open-incentives/internal/domain"
)

//go:generate mockgen -source=main.go -destination=mocks/mocks.go -package=mocks

type ScopedProjectStore interface {
	List(ctx context.Context) ([]domain.Project, error)
}

type Usecase struct {
	projects ScopedProjectStore
}

func New(projects ScopedProjectStore) *Usecase {
	return &Usecase{projects: projects}
}

func (uc *Usecase) Execute(ctx context.Context) (domain.ListProjectsUsecaseOutput, error) {
	projects, err := uc.projects.List(ctx)
	if err != nil {
		return domain.ListProjectsUsecaseOutput{}, err
	}

	return domain.ListProjectsUsecaseOutput{
		Projects: projects,
	}, nil
}
