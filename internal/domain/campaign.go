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
	Slug      string
	ProjectID int64
	OrgId     int64
	Name      string
	Status    CampaignStatus
	Rule      json.RawMessage
}
