package adapters

import (
	"context"
	"encoding/json"
	"fmt"

	engine "github.com/squeakycheese75/open-incentives-engine"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type Adapter struct {
	runtime engine.Runtime
}

func New(runtime engine.Runtime) *Adapter {
	return &Adapter{
		runtime: runtime,
	}
}

func (a *Adapter) Evaluate(
	ctx context.Context,
	req domain.RuleEvaluationRequest,
) (domain.RuleEvaluationResult, error) {
	var rules []engine.Rule

	if err := json.Unmarshal(req.RawRules, &rules); err != nil {
		return domain.RuleEvaluationResult{},
			fmt.Errorf("decode runtime rules: %w", err)
	}

	result, err := a.runtime.Evaluate(ctx, engine.EvaluationRequest{
		Facts: req.Facts,
		Rules: rules,
	})
	if err != nil {
		return domain.RuleEvaluationResult{}, err
	}

	return mapResult(result, rules), nil
}

func (a *Adapter) ValidateRules(rules json.RawMessage) error {
	return a.runtime.ValidateRules(rules)
}
