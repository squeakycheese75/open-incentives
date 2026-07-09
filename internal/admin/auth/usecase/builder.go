package usecase_auth

import (
	"context"
	"time"

	"github.com/squeakycheese75/open-incentives/internal/admin/auth/usecase/auth_login"
	"github.com/squeakycheese75/open-incentives/internal/domain"
)

type (
	AuthLoginUsecase interface {
		Execute(ctx context.Context, input domain.AuthLoginUsecaseInput) (*domain.AuthLoginUsecaseOutput, error)
	}
)

type (
	UserStore interface {
		Find(ctx context.Context, email string, orgID int64) (domain.User, error)
	}
	OrgStore interface {
		Find(ctx context.Context, publicID string) (domain.Organization, error)
	}
	PasswordSvc interface {
		Verify(password, hash string) bool
	}
	TokenSvc interface {
		Create(userID, orgID string) (string, time.Time, error)
	}
)

type AuthUsecaseContainer struct {
	userStore   UserStore
	orgStore    OrgStore
	tokenSvc    TokenSvc
	passwordSvc PasswordSvc
}

func NewUserUsecaseContainer(userStore UserStore, orgStore OrgStore, tokenSvc TokenSvc, passwordSvc PasswordSvc) *AuthUsecaseContainer {
	return &AuthUsecaseContainer{
		userStore:   userStore,
		orgStore:    orgStore,
		tokenSvc:    tokenSvc,
		passwordSvc: passwordSvc,
	}
}

func (uc *AuthUsecaseContainer) LoginUserUsecase() AuthLoginUsecase {
	return auth_login.NewAuthLoginUsecase(uc.userStore, uc.orgStore, uc.passwordSvc, uc.tokenSvc)
}
