package models

import (
	"database/sql"
	"time"
)

type PasswordResetToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	UsedAt    sql.NullTime
	CreatedAt time.Time
	UpdatedAt time.Time
}
