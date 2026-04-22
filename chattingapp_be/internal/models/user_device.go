package models

import (
	"database/sql"
	"time"
)

type UserDevice struct {
	ID           int64
	UserID        int64
	DeviceUUID    string
	DeviceName    string
	DeviceType    string
	Platform      string
	AppVersion    sql.NullString
	OSVersion     sql.NullString
	PushToken     sql.NullString
	IsTrusted     bool
	IsActive      bool
	LastSeenAt    sql.NullTime
	RegisteredAt  time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}