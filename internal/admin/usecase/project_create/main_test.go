package project_create

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/squeakycheese75/open-incentives/internal/admin/usecase/project_create/mocks"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

func TestUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)

	idGenerator := mocks.NewMockPublicIDGenerator(ctrl)
	projects := mocks.NewMockScopedProjectStore(ctrl)

	created := domain.Project{ID: 1, PublicID: "proj_abc123", Name: "My Project"}

	idGenerator.EXPECT().New(prefix).Return("proj_abc123", nil)
	projects.EXPECT().Create(gomock.Any(), domain.Project{PublicID: "proj_abc123", Name: "My Project"}).Return(created, nil)

	uc := New(idGenerator, projects)

	got, err := uc.Execute(context.Background(), domain.CreateProjectUsecaseInput{Name: "My Project"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	want := domain.CreateProjectUsecaseOutput{Project: created}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Execute() = %#v, want %#v", got, want)
	}
}

func TestUsecase_Execute_MissingName(t *testing.T) {
	ctrl := gomock.NewController(t)

	idGenerator := mocks.NewMockPublicIDGenerator(ctrl)
	projects := mocks.NewMockScopedProjectStore(ctrl)

	uc := New(idGenerator, projects)

	_, err := uc.Execute(context.Background(), domain.CreateProjectUsecaseInput{Name: "  "})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrInvalidInput)
	}
}

func TestUsecase_Execute_StoreError(t *testing.T) {
	ctrl := gomock.NewController(t)

	idGenerator := mocks.NewMockPublicIDGenerator(ctrl)
	projects := mocks.NewMockScopedProjectStore(ctrl)

	storeErr := errors.New("store unavailable")

	idGenerator.EXPECT().New(prefix).Return("proj_abc123", nil)
	projects.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domain.Project{}, storeErr)

	uc := New(idGenerator, projects)

	_, err := uc.Execute(context.Background(), domain.CreateProjectUsecaseInput{Name: "My Project"})
	if !errors.Is(err, storeErr) {
		t.Fatalf("Execute() error = %v, want %v", err, storeErr)
	}
}
