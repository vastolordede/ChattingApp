package models

import (
	"database/sql"
	"time"
)

type ConversationMember struct {
	ID                int64
	ConversationID    int64
	UserID            int64
	Role              string
	JoinedAt          time.Time
	LeftAt            sql.NullTime
	IsActive          bool
	Nickname          sql.NullString
	LastReadMessageID sql.NullInt64
	LastReadAt        sql.NullTime
	MuteUntil         sql.NullTime
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
