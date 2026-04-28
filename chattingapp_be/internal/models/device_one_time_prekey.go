package models

import (
	"database/sql"
	"time"
)

type DeviceOneTimePreKey struct {
	ID        int64
	DeviceID  int64
	KeyID     int
	PublicKey string
	Algorithm string
	Version   int
	IsUsed    bool
	UsedAt    sql.NullTime
	CreatedAt time.Time
	ExpiredAt sql.NullTime
}
