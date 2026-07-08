package eval

import engine "github.com/squeakycheese75/open-incentives-engine"

func toEngineRequest(req EvaluateRequest) engine.EvaluationRequest {
	rules := make([]engine.Rule, 0, len(req.Rules))

	for _, rule := range req.Rules {
		actions := make([]engine.Action, 0, len(rule.Actions))
		for _, action := range rule.Actions {
			actions = append(actions, engine.Action{
				Type:   action.Type,
				Params: action.Params,
			})
		}

		rules = append(rules, engine.Rule{
			ID:         rule.ID,
			Name:       rule.Name,
			Conditions: mapCondition(rule.Conditions),
			Actions:    actions,
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

func mapCondition(c ConditionRequest) engine.Condition {
	return engine.Condition{
		All:      mapConditions(c.All),
		Any:      mapConditions(c.Any),
		Fact:     c.Fact,
		Operator: c.Operator,
		Value:    c.Value,
	}
}

func mapConditions(conditions []ConditionRequest) []engine.Condition {
	if len(conditions) == 0 {
		return nil
	}

	out := make([]engine.Condition, 0, len(conditions))
	for _, c := range conditions {
		out = append(out, mapCondition(c))
	}

	return out
}
