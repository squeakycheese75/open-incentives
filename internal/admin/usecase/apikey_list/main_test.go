package apikey_list

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/squeakycheese75/open-incentives/internal/admin/usecase/apikey_list/mocks"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

func TestUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)

	apikeys := mocks.NewMockScopedAPIKeyStore(ctrl)
	projects := mocks.NewMockScopedProjectStore(ctrl)

	project := domain.Project{ID: 7, PublicID: "proj_123"}
	expected := []domain.APIKey{
		{PublicID: "api_1", Name: "key one", Prefix: "api", Status: domain.APIKeyStatusActive},
	}

	projects.EXPECT().Find(gomock.Any(), "proj_123").Return(project, nil)
	apikeys.EXPECT().List(gomock.Any(), int64(7)).Return(expected, nil)

	uc := New(apikeys, projects)

	got, err := uc.Execute(context.Background(), domain.ListAPIKeysUsecaseInput{ProjectPublicID: "proj_123"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	want := domain.ListAPIKeysUsecaseOutput{APIKeys: expected}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Execute() = %#v, want %#v", got, want)
	}
}

func TestUsecase_Execute_MissingProjectID(t *testing.T) {
	ctrl := gomock.NewController(t)

	apikeys := mocks.NewMockScopedAPIKeyStore(ctrl)
	projects := mocks.NewMockScopedProjectStore(ctrl)

	uc := New(apikeys, projects)

	_, err := uc.Execute(context.Background(), domain.ListAPIKeysUsecaseInput{ProjectPublicID: "  "})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrInvalidInput)
	}
}

func TestUsecase_Execute_ProjectNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	apikeys := mocks.NewMockScopedAPIKeyStore(ctrl)
	projects := mocks.NewMockScopedProjectStore(ctrl)

	projects.EXPECT().Find(gomock.Any(), "proj_123").Return(domain.Project{}, domain.ErrNotFound)

	uc := New(apikeys, projects)

	_, err := uc.Execute(context.Background(), domain.ListAPIKeysUsecaseInput{ProjectPublicID: "proj_123"})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrNotFound)
	}
}
