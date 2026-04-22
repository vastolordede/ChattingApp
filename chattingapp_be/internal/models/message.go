package models

import (
	"database/sql"
	"time"
)

type Message struct {
	ID                     int64
	ConversationID         int64
	SenderUserID           int64
	MessageType            string
	Content                sql.NullString
	ReplyToMessageID       sql.NullInt64
	ForwardedFromMessageID sql.NullInt64
	Status                 string
	IsEdited               bool
	EditedAt               sql.NullTime
	IsDeleted              bool
	DeletedAt              sql.NullTime
	ClientMessageID        sql.NullString
	SentAt                 time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}