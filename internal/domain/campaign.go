package domain

import (
	"encoding/json"
)

type CampaignStatus string

const (
	CampaignStatusActive   CampaignStatus = "active"
	CampaignStatusInactive CampaignStatus = "inactive"
)

type Campaign struct {
	PublicID  string
	ProjectID int64
	OrgId     int64
	Name      string
	KeyHash   string
	Status    CampaignStatus
	Rules     json.RawMessage
}
