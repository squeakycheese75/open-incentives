package domain

import "encoding/json"

type AppliedDiscount struct {
	CampaignID   string
	CampaignName string
	RuleID       string
	Type         string
	Amount       float64
}

type EvaluateDecision struct {
	Matched          bool
	CampaignsMatched int
}

type EvaluateCartOutput struct {
	Subtotal      float64
	DiscountTotal float64
	Total         float64
	Currency      string
}

type EvaluateUsecaseOutput struct {
	Decision  EvaluateDecision
	Cart      EvaluateCartOutput
	Discounts []AppliedDiscount
}

type RuleEvaluationRequest struct {
	Facts    map[string]any
	RawRules json.RawMessage
}

type RuleEvaluationResult struct {
	MatchedRules []MatchedRule
}

type MatchedRule struct {
	RuleID  string
	Actions []Action
}

type MatchedAction struct {
	CampaignID   string
	CampaignName string
	RuleID       string
	Action       Action
}

type Action struct {
	Type   string
	Params map[string]any
}
