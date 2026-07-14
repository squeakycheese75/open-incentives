package usecase_eval

import (
	"context"

	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/evaluate/usecase/evaluate"
	"github.com/squeakycheese75/open-incentives/internal/store"
)

type EvaluateUsecase interface {
	Execute(ctx context.Context, input domain.EvaluateUsecaseInput) (domain.EvaluateUsecaseOutput, error)
}

type (
	RuntimeAdapter interface {
		Evaluate(
			ctx context.Context,
			req domain.RuleEvaluationRequest,
		) (domain.RuleEvaluationResult, error)
	}
	ActionApplier interface {
		Apply(
			cart domain.Cart,
			matches []domain.MatchedAction,
		) (domain.EvaluateUsecaseOutput, error)
	}
)

type EvalUsecaseFactory struct {
	runtimeAdapter RuntimeAdapter
	actionApplier  ActionApplier
	campaignStore  store.CampaignStore
}

func NewAdminUsecaseFactory(campaignStore store.CampaignStore, runtimeAdapter RuntimeAdapter, actionApplier ActionApplier) *EvalUsecaseFactory {
	return &EvalUsecaseFactory{
		actionApplier:  actionApplier,
		runtimeAdapter: runtimeAdapter,
		campaignStore:  campaignStore,
	}
}

func (f *EvalUsecaseFactory) EvaluateUsecase(orgID int64) EvaluateUsecase {
	return evaluate.NewEvaluateUseCase(f.campaignStore.Scope(orgID), f.runtimeAdapter, f.actionApplier)
}
