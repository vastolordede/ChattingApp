package models

import (
	"database/sql"
	"time"
)

type DeviceSignedPreKey struct {
	ID        int64
	DeviceID  int64
	KeyID     int
	PublicKey string
	Signature string
	Algorithm string
	Version   int
	IsActive  bool
	CreatedAt time.Time
	ExpiredAt sql.NullTime
	RevokedAt sql.NullTime
}