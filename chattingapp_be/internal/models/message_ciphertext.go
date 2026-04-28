package models

import (
	"database/sql"
	"time"
)

type MessageCiphertext struct {
	ID               int64
	MessageID        int64
	TargetDeviceID   int64
	SenderDeviceID   sql.NullInt64
	Ciphertext       string
	EncryptionHeader sql.NullString
	Nonce            sql.NullString
	Algorithm        string
	MessageVersion   int
	IsDelivered      bool
	DeliveredAt      sql.NullTime
	CreatedAt        time.Time
}
