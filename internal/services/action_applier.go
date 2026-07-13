package services

import (
	"fmt"
	"math"

	"github.com/squeakycheese75/open-incentives/internal/domain"
)

const (
	ActionTypePercentageDiscount = "percentage_discount"
	ActionTypeFixedDiscount      = "fixed_discount"
)

type ActionApplier interface {
	Apply(
		cart domain.Cart,
		matches []domain.MatchedAction,
	) (domain.EvaluateUsecaseOutput, error)
}

type DefaultActionApplier struct{}

func NewActionApplier() *DefaultActionApplier {
	return &DefaultActionApplier{}
}

func (a *DefaultActionApplier) Apply(
	cart domain.Cart,
	matches []domain.MatchedAction,
) (domain.EvaluateUsecaseOutput, error) {
	total := cart.Subtotal
	discounts := make([]domain.AppliedDiscount, 0, len(matches))
	matchedCampaigns := make(map[string]struct{})

	for _, match := range matches {
		amount, err := applyAction(total, match.Action)
		if err != nil {
			return domain.EvaluateUsecaseOutput{}, fmt.Errorf(
				"apply campaign %s rule %s: %w",
				match.CampaignID,
				match.RuleID,
				err,
			)
		}

		if amount <= 0 {
			continue
		}

		if amount > total {
			amount = total
		}

		total = roundMoney(total - amount)

		discounts = append(discounts, domain.AppliedDiscount{
			CampaignID:   match.CampaignID,
			CampaignName: match.CampaignName,
			RuleID:       match.RuleID,
			Type:         match.Action.Type,
			Amount:       amount,
		})

		matchedCampaigns[match.CampaignID] = struct{}{}
	}

	discountTotal := roundMoney(cart.Subtotal - total)

	return domain.EvaluateUsecaseOutput{
		Decision: domain.EvaluateDecision{
			Matched:          len(discounts) > 0,
			CampaignsMatched: len(matchedCampaigns),
		},
		Cart: domain.EvaluateCartOutput{
			Subtotal:      cart.Subtotal,
			DiscountTotal: discountTotal,
			Total:         total,
			Currency:      cart.Currency,
		},
		Discounts: discounts,
	}, nil
}

func applyAction(
	currentTotal float64,
	action domain.Action,
) (float64, error) {
	switch action.Type {
	case ActionTypePercentageDiscount:
		value, err := numericParam(action.Params, "value")
		if err != nil {
			return 0, err
		}

		if value < 0 || value > 100 {
			return 0, fmt.Errorf(
				"percentage discount value must be between 0 and 100",
			)
		}

		return roundMoney(currentTotal * value / 100), nil

	case ActionTypeFixedDiscount:
		value, err := numericParam(action.Params, "value")
		if err != nil {
			return 0, err
		}

		if value < 0 {
			return 0, fmt.Errorf(
				"fixed discount value must not be negative",
			)
		}

		return roundMoney(value), nil

	default:
		return 0, fmt.Errorf(
			"unsupported action type %q",
			action.Type,
		)
	}
}

func numericParam(
	params map[string]any,
	name string,
) (float64, error) {
	value, ok := params[name]
	if !ok {
		return 0, fmt.Errorf("missing action parameter %q", name)
	}

	switch number := value.(type) {
	case float64:
		return number, nil
	case float32:
		return float64(number), nil
	case int:
		return float64(number), nil
	case int8:
		return float64(number), nil
	case int16:
		return float64(number), nil
	case int32:
		return float64(number), nil
	case int64:
		return float64(number), nil
	case uint:
		return float64(number), nil
	case uint8:
		return float64(number), nil
	case uint16:
		return float64(number), nil
	case uint32:
		return float64(number), nil
	case uint64:
		return float64(number), nil
	default:
		return 0, fmt.Errorf(
			"action parameter %q must be numeric",
			name,
		)
	}
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}
