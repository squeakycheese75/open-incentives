package engine

import "context"

type Engine interface {
	Evaluate(ctx context.Context, req EvaluationRequest) (EvaluationResult, error)
}

type EvaluationRequest struct {
	Facts map[string]any
	Rules []Rule
}

type Rule struct {
	ID         string
	Name       string
	Conditions Conditions
	Actions    []Action
}

type Conditions struct {
	All []Condition
}

type Condition struct {
	Fact     string
	Operator string
	Value    any
}

type Action struct {
	Type   string
	Params map[string]any
}

type EvaluationResult struct {
	Actions       []Action
	MatchedRules  []string
	RejectedRules []string
	Trace         []TraceEntry
}

type TraceEntry struct {
	RuleID   string
	Fact     string
	Operator string
	Expected any
	Actual   any
	Passed   bool
	Message  string
}
