package models

import (
	"database/sql"
	"time"
)

type DeviceIdentityKey struct {
	ID          int64
	DeviceID    int64
	PublicKey   string
	Algorithm   string
	Fingerprint string
	Version     int
	IsActive    bool
	CreatedAt   time.Time
	ExpiredAt   sql.NullTime
	RevokedAt   sql.NullTime
}
