package admin

import (
	usecase_admin "github.com/squeakycheese75/open-incentives/internal/admin/usecase"
)

type Handler struct {
	adminContainer *usecase_admin.AdminUsecaseFactory
}

func NewHandler(adminContainer *usecase_admin.AdminUsecaseFactory) *Handler {
	return &Handler{
		adminContainer: adminContainer,
	}
}
