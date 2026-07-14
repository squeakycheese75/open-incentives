// internal/eval/usecase/evaluate/main_test.go
package evaluate

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/evaluate/usecase/evaluate/mocks"
)

func TestUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)

	campaigns := mocks.NewMockScopedCampaignStore(ctrl)
	runtime := mocks.NewMockRuntime(ctrl)
	applier := mocks.NewMockActionApplier(ctrl)

	rawRules := json.RawMessage(`[
		{
			"id": "rule_10_percent_over_50"
		}
	]`)

	input := domain.EvaluateUsecaseInput{
		ProjectID: 20,
		Customer: domain.Customer{
			ID:      "user_123",
			Country: "DE",
			Tier:    "gold",
		},
		Cart: domain.Cart{
			Subtotal: 50,
			Currency: "EUR",
		},
	}

	campaign := domain.Campaign{
		PublicID: "camp_123",
		Name:     "10% off orders over €50",
		Status:   domain.CampaignStatusActive,
		Rules:    rawRules,
	}

	runtimeResult := domain.RuleEvaluationResult{
		MatchedRules: []domain.MatchedRule{
			{
				RuleID: "rule_10_percent_over_50",
				Actions: []domain.Action{
					{
						Type: "percentage_discount",
						Params: map[string]any{
							"value": float64(10),
						},
					},
				},
			},
		},
	}

	expectedMatches := []domain.MatchedAction{
		{
			CampaignID:   "camp_123",
			CampaignName: "10% off orders over €50",
			RuleID:       "rule_10_percent_over_50",
			Action: domain.Action{
				Type: "percentage_discount",
				Params: map[string]any{
					"value": float64(10),
				},
			},
		},
	}

	expectedOutput := domain.EvaluateUsecaseOutput{
		Decision: domain.EvaluateDecision{
			Matched:          true,
			CampaignsMatched: 1,
		},
		Cart: domain.EvaluateCartOutput{
			Subtotal:      50,
			DiscountTotal: 5,
			Total:         45,
			Currency:      "EUR",
		},
	}

	campaigns.EXPECT().
		List(gomock.Any(), int64(20)).
		Return([]domain.Campaign{campaign}, nil)

	runtime.EXPECT().
		Evaluate(gomock.Any(), gomock.Any()).
		DoAndReturn(func(
			_ context.Context,
			req domain.RuleEvaluationRequest,
		) (domain.RuleEvaluationResult, error) {
			expectedFacts := map[string]any{
				"customer.id":      "user_123",
				"customer.country": "DE",
				"customer.tier":    "gold",
				"cart.subtotal":    float64(50),
				"cart.currency":    "EUR",
			}

			if !reflect.DeepEqual(req.Facts, expectedFacts) {
				t.Fatalf("Facts = %#v, want %#v", req.Facts, expectedFacts)
			}

			if !reflect.DeepEqual(req.RawRules, rawRules) {
				t.Fatalf("RawRules = %s, want %s", req.RawRules, rawRules)
			}

			return runtimeResult, nil
		})

	applier.EXPECT().
		Apply(input.Cart, expectedMatches).
		Return(expectedOutput, nil)

	uc := NewEvaluateUseCase(campaigns, runtime, applier)

	got, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !reflect.DeepEqual(got, expectedOutput) {
		t.Fatalf("Execute() = %#v, want %#v", got, expectedOutput)
	}
}

func TestUsecase_Execute_SkipsInactiveCampaigns(t *testing.T) {
	ctrl := gomock.NewController(t)

	campaigns := mocks.NewMockScopedCampaignStore(ctrl)
	runtime := mocks.NewMockRuntime(ctrl)
	applier := mocks.NewMockActionApplier(ctrl)

	input := domain.EvaluateUsecaseInput{
		ProjectID: 20,
		Cart: domain.Cart{
			Subtotal: 50,
			Currency: "EUR",
		},
	}

	campaigns.EXPECT().
		List(gomock.Any(), int64(20)).
		Return([]domain.Campaign{
			{
				PublicID: "camp_inactive",
				Status:   domain.CampaignStatusInactive,
				Rules:    json.RawMessage(`[]`),
			},
		}, nil)

	expected := domain.EvaluateUsecaseOutput{
		Cart: domain.EvaluateCartOutput{
			Subtotal: 50,
			Total:    50,
			Currency: "EUR",
		},
	}

	applier.EXPECT().
		Apply(input.Cart, []domain.MatchedAction(nil)).
		Return(expected, nil)

	uc := NewEvaluateUseCase(campaigns, runtime, applier)

	got, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Execute() = %#v, want %#v", got, expected)
	}
}

func TestUsecase_Execute_ReturnsStoreError(t *testing.T) {
	ctrl := gomock.NewController(t)

	campaigns := mocks.NewMockScopedCampaignStore(ctrl)
	runtime := mocks.NewMockRuntime(ctrl)
	applier := mocks.NewMockActionApplier(ctrl)

	storeErr := errors.New("store unavailable")

	campaigns.EXPECT().
		List(gomock.Any(), int64(20)).
		Return(nil, storeErr)

	uc := NewEvaluateUseCase(campaigns, runtime, applier)

	_, err := uc.Execute(context.Background(), domain.EvaluateUsecaseInput{
		ProjectID: 20,
	})

	if !errors.Is(err, storeErr) {
		t.Fatalf("Execute() error = %v, want %v", err, storeErr)
	}
}

func TestUsecase_Execute_ReturnsRuntimeError(t *testing.T) {
	ctrl := gomock.NewController(t)

	campaigns := mocks.NewMockScopedCampaignStore(ctrl)
	runtime := mocks.NewMockRuntime(ctrl)
	applier := mocks.NewMockActionApplier(ctrl)

	runtimeErr := errors.New("runtime unavailable")

	campaigns.EXPECT().
		List(gomock.Any(), int64(20)).
		Return([]domain.Campaign{
			{
				PublicID: "camp_123",
				Status:   domain.CampaignStatusActive,
				Rules:    json.RawMessage(`[]`),
			},
		}, nil)

	runtime.EXPECT().
		Evaluate(gomock.Any(), gomock.Any()).
		Return(domain.RuleEvaluationResult{}, runtimeErr)

	uc := NewEvaluateUseCase(campaigns, runtime, applier)

	_, err := uc.Execute(context.Background(), domain.EvaluateUsecaseInput{
		ProjectID: 20,
	})

	if !errors.Is(err, runtimeErr) {
		t.Fatalf("Execute() error = %v, want %v", err, runtimeErr)
	}
}

func TestUsecase_Execute_ReturnsActionApplierError(t *testing.T) {
	ctrl := gomock.NewController(t)

	campaigns := mocks.NewMockScopedCampaignStore(ctrl)
	runtime := mocks.NewMockRuntime(ctrl)
	applier := mocks.NewMockActionApplier(ctrl)

	applyErr := errors.New("apply failed")

	input := domain.EvaluateUsecaseInput{
		ProjectID: 20,
		Cart: domain.Cart{
			Subtotal: 50,
			Currency: "EUR",
		},
	}

	campaigns.EXPECT().
		List(gomock.Any(), int64(20)).
		Return(nil, nil)

	applier.EXPECT().
		Apply(input.Cart, []domain.MatchedAction(nil)).
		Return(domain.EvaluateUsecaseOutput{}, applyErr)

	uc := NewEvaluateUseCase(campaigns, runtime, applier)

	_, err := uc.Execute(context.Background(), input)

	if !errors.Is(err, applyErr) {
		t.Fatalf("Execute() error = %v, want %v", err, applyErr)
	}
}
