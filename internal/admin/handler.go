package admin

import (
	usecase_admin "github.com/squeakycheese75/open-incentives/internal/admin/usecase"
)

type Handler struct {
	adminContainer *usecase_admin.AdminUsecaseContainer
}

func NewHandler(adminContainer *usecase_admin.AdminUsecaseContainer) *Handler {
	return &Handler{
		adminContainer: adminContainer,
	}
}
