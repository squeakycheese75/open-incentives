package project_create_apikey

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/store"
)

const prefix string = "api"

type (
	CryptoSvc interface {
		GenerateKey(size int) ([]byte, error)
	}
	PublicIDGenerator interface {
		New(prefix string) (string, error)
	}
	APIKeyStore interface {
		Create(ctx context.Context, in domain.APIKey) (domain.APIKey, error)
	}
	ProjectStore interface {
		Find(ctx context.Context, publicID string) (domain.Project, error)
	}
	PasswordSvc interface {
		Hash(password string) (string, error)
	}
)

type Usecase struct {
	cryptoSvc   CryptoSvc
	idGenerator PublicIDGenerator
	apiKeyStore store.APIKeyStore
	projects    store.ProjectStore
	passwordSvc PasswordSvc
}

func New(cryptoSvc CryptoSvc, idGenerator PublicIDGenerator, apiKeyStore store.APIKeyStore, passwordSvc PasswordSvc, projects store.ProjectStore) *Usecase {
	return &Usecase{
		cryptoSvc:   cryptoSvc,
		apiKeyStore: apiKeyStore,
		idGenerator: idGenerator,
		passwordSvc: passwordSvc,
		projects:    projects,
	}
}

func (uc *Usecase) Execute(ctx context.Context, in domain.CreateProjectAPIKEYUsecaseInput) (domain.CreateProjectAPIKEYUsecaseOutput, error) {
	projectPublicID := strings.TrimSpace(in.ProjectPublicID)
	if projectPublicID == "" {
		return domain.CreateProjectAPIKEYUsecaseOutput{}, fmt.Errorf("project id is required: %w", domain.ErrInvalidInput)
	}

	project, err := uc.projects.Scope(in.OrgID).Find(ctx, projectPublicID)
	if err != nil {
		return domain.CreateProjectAPIKEYUsecaseOutput{}, err
	}

	apikeyPublicID, err := uc.idGenerator.New(prefix)
	if err != nil {
		return domain.CreateProjectAPIKEYUsecaseOutput{}, err
	}

	secretBytes, err := uc.cryptoSvc.GenerateKey(32)
	if err != nil {
		return domain.CreateProjectAPIKEYUsecaseOutput{}, fmt.Errorf("failed to generate api key secret: %v", err)
	}

	secret := base64.RawURLEncoding.EncodeToString(secretBytes)

	hash, err := uc.passwordSvc.Hash(secret)
	if err != nil {
		return domain.CreateProjectAPIKEYUsecaseOutput{}, fmt.Errorf("failed to hash api key secret: %v", err)
	}

	out, err := uc.apiKeyStore.Scope(in.OrgID).Create(ctx, domain.APIKey{
		Name:      in.Name,
		PublicID:  apikeyPublicID,
		ProjectID: project.ID,
		KeyHash:   hash,
		Prefix:    prefix,
		Status:    domain.APIKeyStatusActive,
	})
	if err != nil {
		return domain.CreateProjectAPIKEYUsecaseOutput{}, fmt.Errorf("failed to persist api key: %v", err)
	}
	return domain.CreateProjectAPIKEYUsecaseOutput{
		APIKeyPublicID: out.PublicID,
		APIKey:         fmt.Sprintf("%s.%s", apikeyPublicID, secret),
		CreatedAt:      out.CreatedAt,
	}, nil
}
