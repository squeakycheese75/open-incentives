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
	Slug   string
	Name   string
	Status CampaignStatus
	Rule   json.RawMessage
}
