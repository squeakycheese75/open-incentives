package campaign_update

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/squeakycheese75/open-incentives/internal/admin/usecase/campaign_update/mocks"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

func TestUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)

	campaigns := mocks.NewMockScopedCampaignStore(ctrl)
	projects := mocks.NewMockScopedProjectStore(ctrl)

	rawRules := json.RawMessage(`[{"id":"rule_1"}]`)

	input := domain.UpdateCampaignUsecaseInput{
		CampaignPublicID: "camp_123",
		ProjectPublicID:  "proj_123",
		Name:             "10% off orders over €50",
		Status:           "active",
		Rules:            rawRules,
	}

	project := domain.Project{ID: 7, PublicID: "proj_123"}

	updated := domain.Campaign{
		PublicID:  "camp_123",
		ProjectID: 7,
		Name:      "10% off orders over €50",
		Status:    domain.CampaignStatusActive,
		Rules:     rawRules,
	}

	projects.EXPECT().Find(gomock.Any(), "proj_123").Return(project, nil)
	campaigns.EXPECT().Update(gomock.Any(), "camp_123", int64(7), domain.Campaign{
		Name:   "10% off orders over €50",
		Status: domain.CampaignStatusActive,
		Rules:  rawRules,
	}).Return(updated, nil)

	uc := New(campaigns, projects)

	got, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	want := domain.UpdateCampaignUsecaseOutput{Campaign: updated}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Execute() = %#v, want %#v", got, want)
	}
}

func TestUsecase_Execute_MissingName(t *testing.T) {
	ctrl := gomock.NewController(t)

	campaigns := mocks.NewMockScopedCampaignStore(ctrl)
	projects := mocks.NewMockScopedProjectStore(ctrl)

	uc := New(campaigns, projects)

	_, err := uc.Execute(context.Background(), domain.UpdateCampaignUsecaseInput{
		CampaignPublicID: "camp_123",
		ProjectPublicID:  "proj_123",
		Name:             "  ",
		Status:           "active",
	})

	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrInvalidInput)
	}
}

func TestUsecase_Execute_InvalidStatus(t *testing.T) {
	ctrl := gomock.NewController(t)

	campaigns := mocks.NewMockScopedCampaignStore(ctrl)
	projects := mocks.NewMockScopedProjectStore(ctrl)

	uc := New(campaigns, projects)

	_, err := uc.Execute(context.Background(), domain.UpdateCampaignUsecaseInput{
		CampaignPublicID: "camp_123",
		ProjectPublicID:  "proj_123",
		Name:             "Some Campaign",
		Status:           "bogus",
	})

	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrInvalidInput)
	}
}

func TestUsecase_Execute_ProjectNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	campaigns := mocks.NewMockScopedCampaignStore(ctrl)
	projects := mocks.NewMockScopedProjectStore(ctrl)

	projects.EXPECT().Find(gomock.Any(), "proj_123").Return(domain.Project{}, domain.ErrNotFound)

	uc := New(campaigns, projects)

	_, err := uc.Execute(context.Background(), domain.UpdateCampaignUsecaseInput{
		CampaignPublicID: "camp_123",
		ProjectPublicID:  "proj_123",
		Name:             "Some Campaign",
		Status:           "active",
	})

	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestUsecase_Execute_CampaignNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	campaigns := mocks.NewMockScopedCampaignStore(ctrl)
	projects := mocks.NewMockScopedProjectStore(ctrl)

	project := domain.Project{ID: 7, PublicID: "proj_123"}

	projects.EXPECT().Find(gomock.Any(), "proj_123").Return(project, nil)
	campaigns.EXPECT().Update(gomock.Any(), "camp_123", int64(7), gomock.Any()).Return(domain.Campaign{}, domain.ErrNotFound)

	uc := New(campaigns, projects)

	_, err := uc.Execute(context.Background(), domain.UpdateCampaignUsecaseInput{
		CampaignPublicID: "camp_123",
		ProjectPublicID:  "proj_123",
		Name:             "Some Campaign",
		Status:           "active",
	})

	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrNotFound)
	}
}
