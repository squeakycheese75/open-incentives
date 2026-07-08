package admin

type CreateCampaignRequest struct {
	Name   string        `json:"name"`
	Status string        `json:"status"`
	Rules  []RuleRequest `json:"rules"`
}

type RuleRequest struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Conditions ConditionRequest `json:"conditions"`
	Actions    []ActionRequest  `json:"actions"`
}

type ConditionRequest struct {
	All      []ConditionRequest `json:"all,omitempty"`
	Any      []ConditionRequest `json:"any,omitempty"`
	Fact     string             `json:"fact,omitempty"`
	Operator string             `json:"operator,omitempty"`
	Value    any                `json:"value,omitempty"`
}

type ActionRequest struct {
	Type   string         `json:"type"`
	Params map[string]any `json:"params"`
}
