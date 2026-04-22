package models

import (
	"database/sql"
	"time"
)

type Conversation struct {
	ID               int64
	ConversationType string
	Title            sql.NullString
	AvatarURL        sql.NullString
	CreatedByUserID  sql.NullInt64
	Status           string
	LastMessageID    sql.NullInt64
	LastMessageAt    sql.NullTime
	CreatedAt        time.Time
	UpdatedAt        time.Time
}