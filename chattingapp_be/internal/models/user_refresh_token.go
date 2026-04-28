package models

import (
	"database/sql"
	"time"
)

type UserRefreshToken struct {
	ID           int64
	UserID       int64
	UserDeviceID sql.NullInt64
	TokenHash    string
	ExpiresAt    time.Time
	RevokedAt    sql.NullTime
	LastUsedAt   sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
