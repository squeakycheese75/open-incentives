package engine

type EvaluationRequest struct {
	Facts map[string]any
	Rules []Rule
}

type Rule struct {
	ID         string
	Version    int
	Name       string
	Conditions Condition
	Actions    []Action
}

type Condition struct {
	All      []Condition
	Any      []Condition
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
