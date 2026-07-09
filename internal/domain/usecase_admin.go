package domain

type CreateCampaignUsecaseInput struct {
	OrgID           int64
	ProjectPublicID string
	Name            string
	Rules           any
}

type CreateCampaignUsecaseOutput struct {
	CampaignID       int64
	CampaignPublicID string
}
