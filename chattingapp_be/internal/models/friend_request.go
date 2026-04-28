package models

import (
	"database/sql"
	"time"
)

type FriendRequest struct {
	ID          int64
	Status      string
	Message     sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
	RespondedAt sql.NullTime
	ExpiredAt   sql.NullTime
}
