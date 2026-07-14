package project_delete

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/squeakycheese75/open-incentives/internal/admin/usecase/project_delete/mocks"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

func TestUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)

	projects := mocks.NewMockScopedProjectStore(ctrl)

	projects.EXPECT().Delete(gomock.Any(), "proj_1").Return(nil)

	uc := New(projects)

	err := uc.Execute(context.Background(), domain.DeleteProjectUsecaseInput{ProjectPublicID: "proj_1"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestUsecase_Execute_MissingProjectID(t *testing.T) {
	ctrl := gomock.NewController(t)

	projects := mocks.NewMockScopedProjectStore(ctrl)

	uc := New(projects)

	err := uc.Execute(context.Background(), domain.DeleteProjectUsecaseInput{ProjectPublicID: "  "})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrInvalidInput)
	}
}

func TestUsecase_Execute_StoreError(t *testing.T) {
	ctrl := gomock.NewController(t)

	projects := mocks.NewMockScopedProjectStore(ctrl)

	storeErr := errors.New("store unavailable")

	projects.EXPECT().Delete(gomock.Any(), "proj_1").Return(storeErr)

	uc := New(projects)

	err := uc.Execute(context.Background(), domain.DeleteProjectUsecaseInput{ProjectPublicID: "proj_1"})
	if !errors.Is(err, storeErr) {
		t.Fatalf("Execute() error = %v, want %v", err, storeErr)
	}
}
