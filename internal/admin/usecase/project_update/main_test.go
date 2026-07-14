package project_update

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/squeakycheese75/open-incentives/internal/admin/usecase/project_update/mocks"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

func TestUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)

	projects := mocks.NewMockScopedProjectStore(ctrl)

	updated := domain.Project{ID: 1, PublicID: "proj_1", Name: "Renamed"}

	projects.EXPECT().Update(gomock.Any(), "proj_1", "Renamed").Return(updated, nil)

	uc := New(projects)

	got, err := uc.Execute(context.Background(), domain.UpdateProjectUsecaseInput{
		ProjectPublicID: "proj_1",
		Name:            "Renamed",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	want := domain.UpdateProjectUsecaseOutput{Project: updated}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Execute() = %#v, want %#v", got, want)
	}
}

func TestUsecase_Execute_MissingFields(t *testing.T) {
	ctrl := gomock.NewController(t)

	projects := mocks.NewMockScopedProjectStore(ctrl)

	uc := New(projects)

	_, err := uc.Execute(context.Background(), domain.UpdateProjectUsecaseInput{
		ProjectPublicID: "",
		Name:            "Renamed",
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrInvalidInput)
	}
}

func TestUsecase_Execute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	projects := mocks.NewMockScopedProjectStore(ctrl)

	projects.EXPECT().Update(gomock.Any(), "proj_1", "Renamed").Return(domain.Project{}, domain.ErrNotFound)

	uc := New(projects)

	_, err := uc.Execute(context.Background(), domain.UpdateProjectUsecaseInput{
		ProjectPublicID: "proj_1",
		Name:            "Renamed",
	})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrNotFound)
	}
}
