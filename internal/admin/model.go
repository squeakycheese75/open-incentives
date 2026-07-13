package admin

import "encoding/json"

type CreateCampaignRequest struct {
	Name   string          `json:"name"`
	Status string          `json:"status"`
	Rules  json.RawMessage `json:"rules"`
}
