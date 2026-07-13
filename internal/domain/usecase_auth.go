package domain

import "time"

type AuthLoginUsecaseInput struct {
	Email       string
	Password    string
	OrgPublicID string
}

type User struct {
	ID             int64
	PublicID       string
	Email          string
	Role           string
	Status         string
	HashedPassword string
}

type AuthLoginUsecaseOutput struct {
	Token     string
	ExpiresAt time.Time
}
