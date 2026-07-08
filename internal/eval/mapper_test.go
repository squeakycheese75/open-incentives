package eval

import "testing"

func TestToEngineRequestMapsNestedConditions(t *testing.T) {
	req := EvaluateRequest{
		Facts: map[string]any{
			"cart.total":    72,
			"customer.tier": "gold",
		},
		Rules: []RuleRequest{
			{
				ID:   "rule_gold_or_platinum",
				Name: "Gold or platinum discount",
				Conditions: ConditionRequest{
					All: []ConditionRequest{
						{Fact: "cart.total", Operator: "gte", Value: 50},
						{
							Any: []ConditionRequest{
								{Fact: "customer.tier", Operator: "eq", Value: "gold"},
								{Fact: "customer.tier", Operator: "eq", Value: "platinum"},
							},
						},
					},
				},
			},
		},
	}

	got := toEngineRequest(req)

	condition := got.Rules[0].Conditions

	if len(condition.All) != 2 {
		t.Fatalf("len(condition.All) = %d, want 2", len(condition.All))
	}

	if condition.All[0].Fact != "cart.total" {
		t.Fatalf("condition.All[0].Fact = %q, want cart.total", condition.All[0].Fact)
	}

	if len(condition.All[1].Any) != 2 {
		t.Fatalf("len(condition.All[1].Any) = %d, want 2", len(condition.All[1].Any))
	}

	if condition.All[1].Any[0].Value != "gold" {
		t.Fatalf("condition.All[1].Any[0].Value = %v, want gold", condition.All[1].Any[0].Value)
	}
}
