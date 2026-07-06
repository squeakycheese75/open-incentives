package handlers

import (
	"encoding/json"
	"net/http"

	engine "github.com/squeakycheese75/open-incentives-engine"
)

type EvaluateRequest struct {
	Facts map[string]any `json:"facts"`
	Rules []RuleRequest  `json:"rules"`
}

type RuleRequest struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Conditions ConditionsRequest `json:"conditions"`
	Actions    []ActionRequest   `json:"actions"`
}

type ConditionsRequest struct {
	All []ConditionRequest `json:"all"`
}

type ConditionRequest struct {
	Fact     string `json:"fact"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
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

func (s *handlers) Evaluate(w http.ResponseWriter, r *http.Request) {
	var httpReq EvaluateRequest

	if err := json.NewDecoder(r.Body).Decode(&httpReq); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid_json",
		})
		return
	}

	engineReq := toEngineRequest(httpReq)

	result, err := s.engine.Evaluate(r.Context(), engineReq)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "evaluation_failed",
		})
		return
	}

	writeJSON(w, http.StatusOK, toEvaluateResponse(result))
}

func toEngineRequest(req EvaluateRequest) engine.EvaluationRequest {
	rules := make([]engine.Rule, 0, len(req.Rules))

	for _, rule := range req.Rules {
		conditions := make([]engine.Condition, 0, len(rule.Conditions.All))
		for _, condition := range rule.Conditions.All {
			conditions = append(conditions, engine.Condition{
				Fact:     condition.Fact,
				Operator: condition.Operator,
				Value:    condition.Value,
			})
		}

		actions := make([]engine.Action, 0, len(rule.Actions))
		for _, action := range rule.Actions {
			actions = append(actions, engine.Action{
				Type:   action.Type,
				Params: action.Params,
			})
		}

		rules = append(rules, engine.Rule{
			ID:   rule.ID,
			Name: rule.Name,
			Conditions: engine.Conditions{
				All: conditions,
			},
			Actions: actions,
		})
	}

	return engine.EvaluationRequest{
		Facts: req.Facts,
		Rules: rules,
	}
}

func toEvaluateResponse(result engine.EvaluationResult) EvaluateResponse {
	actions := make([]ActionResponse, 0, len(result.Actions))
	for _, action := range result.Actions {
		actions = append(actions, ActionResponse{
			Type:   action.Type,
			Params: action.Params,
		})
	}

	trace := make([]TraceResponse, 0, len(result.Trace))
	for _, entry := range result.Trace {
		trace = append(trace, TraceResponse{
			RuleID:   entry.RuleID,
			Fact:     entry.Fact,
			Operator: entry.Operator,
			Expected: entry.Expected,
			Actual:   entry.Actual,
			Passed:   entry.Passed,
			Message:  entry.Message,
		})
	}

	return EvaluateResponse{
		Actions:       actions,
		MatchedRules:  result.MatchedRules,
		RejectedRules: result.RejectedRules,
		Trace:         trace,
	}
}
