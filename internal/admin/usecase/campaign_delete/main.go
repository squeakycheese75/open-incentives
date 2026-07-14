package campaign_delete

import (
	"context"
	"fmt"
	"strings"

	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type (
	ScopedProjectStore interface {
		Find(ctx context.Context, publicID string) (domain.Project, error)
	}
	ScopedCampaignStore interface {
		Delete(ctx context.Context, campaignPublicID string, projectID int64) error
	}
)

type Usecase struct {
	campaigns ScopedCampaignStore
	projects  ScopedProjectStore
}

func New(campaigns ScopedCampaignStore, projects ScopedProjectStore) *Usecase {
	return &Usecase{campaigns: campaigns, projects: projects}
}

func (uc *Usecase) Execute(ctx context.Context, input domain.DeleteCampaignUsecaseInput) error {
	campaignPublicID := strings.TrimSpace(input.CampaignPublicID)
	if campaignPublicID == "" {
		return fmt.Errorf("campaign id is required: %w", domain.ErrInvalidInput)
	}

	projectPublicID := strings.TrimSpace(input.ProjectPublicID)
	if projectPublicID == "" {
		return fmt.Errorf("project id is required: %w", domain.ErrInvalidInput)
	}

	project, err := uc.projects.Find(ctx, projectPublicID)
	if err != nil {
		return err
	}

	err = uc.campaigns.Delete(ctx, campaignPublicID, project.ID)
	if err != nil {
		return err
	}

	return nil
}
