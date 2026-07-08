package eval

type EvaluateRequest struct {
	Facts map[string]any `json:"facts"`
	// Rules to be removed.
	Rules []RuleRequest `json:"rules"`
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

type EvaluateResponse struct {
	Actions       []ActionResponse `json:"actions"`
	MatchedRules  []string         `json:"matched_rules"`
	RejectedRules []string         `json:"rejected_rules"`
	Trace         []TraceResponse  `json:"trace"`
}

type ActionResponse struct {
	Type   string         `json:"type"`
	Params map[string]any `json:"params"`
}

type TraceResponse struct {
	RuleID   string `json:"rule_id"`
	Fact     string `json:"fact"`
	Operator string `json:"operator"`
	Expected any    `json:"expected"`
	Actual   any    `json:"actual"`
	Passed   bool   `json:"passed"`
	Message  string `json:"message"`
}
