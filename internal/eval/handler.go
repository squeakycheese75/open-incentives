package eval

import (
	usecase_eval "github.com/squeakycheese75/open-incentives/internal/eval/usecase"
)

type Handler struct {
	evalContainer *usecase_eval.EvalUsecaseFactory
}

func NewHandler(evalContainer *usecase_eval.EvalUsecaseFactory) *Handler {
	return &Handler{
		evalContainer: evalContainer,
	}
}
