package domain

import "time"

type Project struct {
	ID        int64
	PublicID  string
	OrgID     int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
