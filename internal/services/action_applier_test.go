package services

import (
	"strings"
	"testing"

	"github.com/squeakycheese75/open-incentives/internal/domain"
)

func TestDefaultActionApplier_Apply_PercentageDiscount(t *testing.T) {
	applier := NewActionApplier()

	output, err := applier.Apply(
		domain.Cart{
			Subtotal: 50,
			Currency: "EUR",
		},
		[]domain.MatchedAction{
			{
				CampaignID:   "camp_123",
				CampaignName: "10% off",
				RuleID:       "rule_10_percent",
				Action: domain.Action{
					Type: ActionTypePercentageDiscount,
					Params: map[string]any{
						"value": float64(10),
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	if !output.Decision.Matched {
		t.Error("Matched = false, want true")
	}

	if output.Decision.CampaignsMatched != 1 {
		t.Errorf(
			"CampaignsMatched = %d, want 1",
			output.Decision.CampaignsMatched,
		)
	}

	if output.Cart.DiscountTotal != 5 {
		t.Errorf(
			"DiscountTotal = %v, want 5",
			output.Cart.DiscountTotal,
		)
	}

	if output.Cart.Total != 45 {
		t.Errorf("Total = %v, want 45", output.Cart.Total)
	}

	if len(output.Discounts) != 1 {
		t.Fatalf(
			"len(Discounts) = %d, want 1",
			len(output.Discounts),
		)
	}

	discount := output.Discounts[0]

	if discount.Amount != 5 {
		t.Errorf("discount Amount = %v, want 5", discount.Amount)
	}

	if discount.Type != ActionTypePercentageDiscount {
		t.Errorf(
			"discount Type = %q, want %q",
			discount.Type,
			ActionTypePercentageDiscount,
		)
	}
}

func TestDefaultActionApplier_Apply_FixedDiscount(t *testing.T) {
	applier := NewActionApplier()

	output, err := applier.Apply(
		domain.Cart{
			Subtotal: 50,
			Currency: "EUR",
		},
		[]domain.MatchedAction{
			{
				CampaignID: "camp_123",
				RuleID:     "rule_fixed_5",
				Action: domain.Action{
					Type: ActionTypeFixedDiscount,
					Params: map[string]any{
						"value": float64(5),
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	if output.Cart.DiscountTotal != 5 {
		t.Errorf(
			"DiscountTotal = %v, want 5",
			output.Cart.DiscountTotal,
		)
	}

	if output.Cart.Total != 45 {
		t.Errorf("Total = %v, want 45", output.Cart.Total)
	}
}

func TestDefaultActionApplier_Apply_StacksActionsInOrder(t *testing.T) {
	applier := NewActionApplier()

	output, err := applier.Apply(
		domain.Cart{
			Subtotal: 100,
			Currency: "EUR",
		},
		[]domain.MatchedAction{
			{
				CampaignID: "camp_percentage",
				RuleID:     "rule_percentage",
				Action: domain.Action{
					Type: ActionTypePercentageDiscount,
					Params: map[string]any{
						"value": float64(10),
					},
				},
			},
			{
				CampaignID: "camp_fixed",
				RuleID:     "rule_fixed",
				Action: domain.Action{
					Type: ActionTypeFixedDiscount,
					Params: map[string]any{
						"value": float64(5),
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	if output.Cart.DiscountTotal != 15 {
		t.Errorf(
			"DiscountTotal = %v, want 15",
			output.Cart.DiscountTotal,
		)
	}

	if output.Cart.Total != 85 {
		t.Errorf("Total = %v, want 85", output.Cart.Total)
	}

	if output.Decision.CampaignsMatched != 2 {
		t.Errorf(
			"CampaignsMatched = %d, want 2",
			output.Decision.CampaignsMatched,
		)
	}
}

func TestDefaultActionApplier_Apply_PercentageDiscountsUseRunningTotal(t *testing.T) {
	applier := NewActionApplier()

	output, err := applier.Apply(
		domain.Cart{
			Subtotal: 100,
			Currency: "EUR",
		},
		[]domain.MatchedAction{
			{
				CampaignID: "camp_1",
				RuleID:     "rule_1",
				Action: domain.Action{
					Type: ActionTypePercentageDiscount,
					Params: map[string]any{
						"value": float64(10),
					},
				},
			},
			{
				CampaignID: "camp_2",
				RuleID:     "rule_2",
				Action: domain.Action{
					Type: ActionTypePercentageDiscount,
					Params: map[string]any{
						"value": float64(10),
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	// 10% of 100 = 10, then 10% of 90 = 9.
	if output.Cart.DiscountTotal != 19 {
		t.Errorf(
			"DiscountTotal = %v, want 19",
			output.Cart.DiscountTotal,
		)
	}

	if output.Cart.Total != 81 {
		t.Errorf("Total = %v, want 81", output.Cart.Total)
	}
}

func TestDefaultActionApplier_Apply_CapsDiscountAtCurrentTotal(t *testing.T) {
	applier := NewActionApplier()

	output, err := applier.Apply(
		domain.Cart{
			Subtotal: 20,
			Currency: "EUR",
		},
		[]domain.MatchedAction{
			{
				CampaignID: "camp_123",
				RuleID:     "rule_fixed_50",
				Action: domain.Action{
					Type: ActionTypeFixedDiscount,
					Params: map[string]any{
						"value": float64(50),
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	if output.Cart.DiscountTotal != 20 {
		t.Errorf(
			"DiscountTotal = %v, want 20",
			output.Cart.DiscountTotal,
		)
	}

	if output.Cart.Total != 0 {
		t.Errorf("Total = %v, want 0", output.Cart.Total)
	}

	if output.Discounts[0].Amount != 20 {
		t.Errorf(
			"discount Amount = %v, want 20",
			output.Discounts[0].Amount,
		)
	}
}

func TestDefaultActionApplier_Apply_NoMatches(t *testing.T) {
	applier := NewActionApplier()

	output, err := applier.Apply(
		domain.Cart{
			Subtotal: 50,
			Currency: "EUR",
		},
		nil,
	)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	if output.Decision.Matched {
		t.Error("Matched = true, want false")
	}

	if output.Decision.CampaignsMatched != 0 {
		t.Errorf(
			"CampaignsMatched = %d, want 0",
			output.Decision.CampaignsMatched,
		)
	}

	if output.Cart.Total != 50 {
		t.Errorf("Total = %v, want 50", output.Cart.Total)
	}

	if len(output.Discounts) != 0 {
		t.Errorf(
			"len(Discounts) = %d, want 0",
			len(output.Discounts),
		)
	}
}

func TestDefaultActionApplier_Apply_ZeroDiscountIsIgnored(t *testing.T) {
	applier := NewActionApplier()

	output, err := applier.Apply(
		domain.Cart{
			Subtotal: 50,
			Currency: "EUR",
		},
		[]domain.MatchedAction{
			{
				CampaignID: "camp_123",
				RuleID:     "rule_zero",
				Action: domain.Action{
					Type: ActionTypePercentageDiscount,
					Params: map[string]any{
						"value": float64(0),
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	if output.Decision.Matched {
		t.Error("Matched = true, want false")
	}

	if len(output.Discounts) != 0 {
		t.Errorf(
			"len(Discounts) = %d, want 0",
			len(output.Discounts),
		)
	}
}

func TestDefaultActionApplier_Apply_UnsupportedAction(t *testing.T) {
	applier := NewActionApplier()

	_, err := applier.Apply(
		domain.Cart{
			Subtotal: 50,
			Currency: "EUR",
		},
		[]domain.MatchedAction{
			{
				CampaignID: "camp_123",
				RuleID:     "rule_unknown",
				Action: domain.Action{
					Type: "unknown_action",
				},
			},
		},
	)

	if err == nil {
		t.Fatal("Apply() error = nil, want error")
	}

	if !strings.Contains(err.Error(), `unsupported action type "unknown_action"`) {
		t.Errorf("Apply() error = %q", err)
	}

	if !strings.Contains(err.Error(), "camp_123") {
		t.Errorf("Apply() error = %q, want campaign context", err)
	}

	if !strings.Contains(err.Error(), "rule_unknown") {
		t.Errorf("Apply() error = %q, want rule context", err)
	}
}

func TestDefaultActionApplier_Apply_MissingValue(t *testing.T) {
	applier := NewActionApplier()

	_, err := applier.Apply(
		domain.Cart{
			Subtotal: 50,
			Currency: "EUR",
		},
		[]domain.MatchedAction{
			{
				CampaignID: "camp_123",
				RuleID:     "rule_missing_value",
				Action: domain.Action{
					Type:   ActionTypePercentageDiscount,
					Params: map[string]any{},
				},
			},
		},
	)

	if err == nil {
		t.Fatal("Apply() error = nil, want error")
	}

	if !strings.Contains(err.Error(), `missing action parameter "value"`) {
		t.Errorf("Apply() error = %q", err)
	}
}

func TestDefaultActionApplier_Apply_NonNumericValue(t *testing.T) {
	applier := NewActionApplier()

	_, err := applier.Apply(
		domain.Cart{
			Subtotal: 50,
			Currency: "EUR",
		},
		[]domain.MatchedAction{
			{
				CampaignID: "camp_123",
				RuleID:     "rule_invalid_value",
				Action: domain.Action{
					Type: ActionTypeFixedDiscount,
					Params: map[string]any{
						"value": "five",
					},
				},
			},
		},
	)

	if err == nil {
		t.Fatal("Apply() error = nil, want error")
	}

	if !strings.Contains(err.Error(), `action parameter "value" must be numeric`) {
		t.Errorf("Apply() error = %q", err)
	}
}

func TestDefaultActionApplier_Apply_InvalidPercentage(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{
			name:  "negative",
			value: -1,
		},
		{
			name:  "over one hundred",
			value: 101,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applier := NewActionApplier()

			_, err := applier.Apply(
				domain.Cart{
					Subtotal: 50,
					Currency: "EUR",
				},
				[]domain.MatchedAction{
					{
						CampaignID: "camp_123",
						RuleID:     "rule_invalid_percentage",
						Action: domain.Action{
							Type: ActionTypePercentageDiscount,
							Params: map[string]any{
								"value": tt.value,
							},
						},
					},
				},
			)

			if err == nil {
				t.Fatal("Apply() error = nil, want error")
			}

			if !strings.Contains(
				err.Error(),
				"percentage discount value must be between 0 and 100",
			) {
				t.Errorf("Apply() error = %q", err)
			}
		})
	}
}

func TestDefaultActionApplier_Apply_NegativeFixedDiscount(t *testing.T) {
	applier := NewActionApplier()

	_, err := applier.Apply(
		domain.Cart{
			Subtotal: 50,
			Currency: "EUR",
		},
		[]domain.MatchedAction{
			{
				CampaignID: "camp_123",
				RuleID:     "rule_negative_fixed",
				Action: domain.Action{
					Type: ActionTypeFixedDiscount,
					Params: map[string]any{
						"value": float64(-5),
					},
				},
			},
		},
	)

	if err == nil {
		t.Fatal("Apply() error = nil, want error")
	}

	if !strings.Contains(
		err.Error(),
		"fixed discount value must not be negative",
	) {
		t.Errorf("Apply() error = %q", err)
	}
}

func TestDefaultActionApplier_Apply_CountsUniqueCampaigns(t *testing.T) {
	applier := NewActionApplier()

	output, err := applier.Apply(
		domain.Cart{
			Subtotal: 100,
			Currency: "EUR",
		},
		[]domain.MatchedAction{
			{
				CampaignID: "camp_123",
				RuleID:     "rule_1",
				Action: domain.Action{
					Type: ActionTypeFixedDiscount,
					Params: map[string]any{
						"value": float64(5),
					},
				},
			},
			{
				CampaignID: "camp_123",
				RuleID:     "rule_2",
				Action: domain.Action{
					Type: ActionTypeFixedDiscount,
					Params: map[string]any{
						"value": float64(5),
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	if output.Decision.CampaignsMatched != 1 {
		t.Errorf(
			"CampaignsMatched = %d, want 1",
			output.Decision.CampaignsMatched,
		)
	}

	if len(output.Discounts) != 2 {
		t.Errorf(
			"len(Discounts) = %d, want 2",
			len(output.Discounts),
		)
	}
}

func TestRoundMoney(t *testing.T) {
	tests := []struct {
		input float64
		want  float64
	}{
		{input: 5.555, want: 5.56},
		{input: 5.554, want: 5.55},
		{input: 0, want: 0},
	}

	for _, tt := range tests {
		got := roundMoney(tt.input)

		if got != tt.want {
			t.Errorf(
				"roundMoney(%v) = %v, want %v",
				tt.input,
				got,
				tt.want,
			)
		}
	}
}
