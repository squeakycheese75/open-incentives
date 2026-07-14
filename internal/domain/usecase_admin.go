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

type DeleteCampaignUsecaseInput struct {
	CampaignPublicID string
	ProjectPublicID  string
}

type UpdateCampaignUsecaseInput struct {
	CampaignPublicID string
	ProjectPublicID  string
	Name             string
	Status           string
	Rules            json.RawMessage
}

type UpdateCampaignUsecaseOutput struct {
	Campaign
}

type ListAPIKeysUsecaseInput struct {
	ProjectPublicID string
}

type ListAPIKeysUsecaseOutput struct {
	APIKeys []APIKey
}

type RevokeAPIKeyUsecaseInput struct {
	APIKeyPublicID  string
	ProjectPublicID string
}

type RevokeAPIKeyUsecaseOutput struct {
	APIKey
}

type CreateProjectUsecaseInput struct {
	Name string
}

type CreateProjectUsecaseOutput struct {
	Project
}

type ListProjectsUsecaseOutput struct {
	Projects []Project
}

type UpdateProjectUsecaseInput struct {
	ProjectPublicID string
	Name            string
}

type UpdateProjectUsecaseOutput struct {
	Project
}

type DeleteProjectUsecaseInput struct {
	ProjectPublicID string
}
