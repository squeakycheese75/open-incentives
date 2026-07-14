package project_list

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/squeakycheese75/open-incentives/internal/admin/usecase/project_list/mocks"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

func TestUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)

	projects := mocks.NewMockScopedProjectStore(ctrl)

	expected := []domain.Project{
		{ID: 1, PublicID: "proj_1", Name: "First"},
		{ID: 2, PublicID: "proj_2", Name: "Second"},
	}

	projects.EXPECT().List(gomock.Any()).Return(expected, nil)

	uc := New(projects)

	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	want := domain.ListProjectsUsecaseOutput{Projects: expected}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Execute() = %#v, want %#v", got, want)
	}
}

func TestUsecase_Execute_StoreError(t *testing.T) {
	ctrl := gomock.NewController(t)

	projects := mocks.NewMockScopedProjectStore(ctrl)

	storeErr := errors.New("store unavailable")

	projects.EXPECT().List(gomock.Any()).Return(nil, storeErr)

	uc := New(projects)

	_, err := uc.Execute(context.Background())
	if !errors.Is(err, storeErr) {
		t.Fatalf("Execute() error = %v, want %v", err, storeErr)
	}
}
