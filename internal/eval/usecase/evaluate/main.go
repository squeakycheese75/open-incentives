package evaluate

import (
	"context"

	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/store"
)

type (
	ActionApplier interface {
		Apply(
			cart domain.Cart,
			matches []domain.MatchedAction,
		) (domain.EvaluateUsecaseOutput, error)
	}
	Runtime interface {
		Evaluate(
			ctx context.Context,
			req domain.RuleEvaluationRequest,
		) (domain.RuleEvaluationResult, error)
	}
	// CampaignStore interface {
	// 	List(ctx context.Context, projectID int64) ([]domain.Campaign, error)
	// }
)

type Usecase struct {
	campaigns     store.CampaignStore
	actionApplier ActionApplier
	runtime       Runtime
}
type EvaluationResult struct{}

func NewEvaluateUseCase(campaigns store.CampaignStore, runtime Runtime, actionApplier ActionApplier) *Usecase {
	return &Usecase{
		campaigns:     campaigns,
		runtime:       runtime,
		actionApplier: actionApplier,
	}
}

func (uc *Usecase) Execute(
	ctx context.Context,
	input domain.EvaluateUsecaseInput,
) (domain.EvaluateUsecaseOutput, error) {
	campaigns, err := uc.campaigns.
		Scope(input.OrgID).
		List(ctx, input.ProjectID)
	if err != nil {
		return domain.EvaluateUsecaseOutput{}, err
	}

	facts := buildFacts(input)

	var matches []domain.MatchedAction

	for _, campaign := range campaigns {
		if campaign.Status != domain.CampaignStatusActive {
			continue
		}

		result, err := uc.runtime.Evaluate(
			ctx,
			domain.RuleEvaluationRequest{
				Facts:    facts,
				RawRules: campaign.Rules,
			},
		)
		if err != nil {
			return domain.EvaluateUsecaseOutput{}, err
		}

		for _, rule := range result.MatchedRules {
			for _, action := range rule.Actions {
				matches = append(matches, domain.MatchedAction{
					CampaignID:   campaign.PublicID,
					CampaignName: campaign.Name,
					RuleID:       rule.RuleID,
					Action:       action,
				})
			}
		}
	}

	return uc.actionApplier.Apply(input.Cart, matches)
}

func buildFacts(input domain.EvaluateUsecaseInput) map[string]any {
	return map[string]any{
		"customer.id":      input.Customer.ID,
		"customer.country": input.Customer.Country,
		"customer.tier":    input.Customer.Tier,
		"cart.subtotal":    input.Cart.Subtotal,
		"cart.currency":    input.Cart.Currency,
	}
}
