package domain

import "time"

type APIKeyStatus string

const (
	APIKeyStatusActive   APIKeyStatus = "active"
	APIKeyStatusInactive APIKeyStatus = "inactive"
)

type APIKey struct {
	Name      string
	PublicID  string
	OrgID     int64
	ProjectID int64
	KeyHash   string
	Prefix    string
	Status    APIKeyStatus
	CreatedAt time.Time
}
