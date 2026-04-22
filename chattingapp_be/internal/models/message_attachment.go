package models

import (
	"database/sql"
	"time"
)

type MessageAttachment struct {
	ID                int64
	MessageID         int64
	AttachmentType    string
	FileName          string
	MimeType          string
	FileSize          int64
	FileURL           string
	ThumbnailURL      sql.NullString
	Width             sql.NullInt64
	Height            sql.NullInt64
	DurationSeconds   sql.NullInt64
	Checksum          sql.NullString
	EncryptionKeyHint sql.NullString
	CreatedAt         time.Time
}