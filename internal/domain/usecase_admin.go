package domain

import (
	"encoding/json"
	"time"
)

type CreateCampaignUsecaseInput struct {
	ProjectPublicID string
	Name            string
	Rules           json.RawMessage
}

type CreateCampaignUsecaseOutput struct {
	CampaignID       int64
	CampaignPublicID string
}

type GetCampaignUsecaseInput struct {
	CampaignPublicID string
	ProjectPublicID  string
}

type GetCampaignUsecaseOutput struct {
	Campaign
}

type CreateProjectAPIKEYUsecaseInput struct {
	OrgID           int64
	ProjectPublicID string
	Name            string
	Description     string
}

type CreateProjectAPIKEYUsecaseOutput struct {
	APIKeyPublicID string
	APIKey         string
	CreatedAt      time.Time
}

type ListCampaignsUsecaseInput struct {
	ProjectPublicID string
}

type ListCampaignsUsecaseOutput struct {
	Campaigns []Campaign
}
