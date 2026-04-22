package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64
	Username     string
	FullName     string
	Email        string
	PhoneNumber  string
	PasswordHash string
	AvatarURL    sql.NullString
	Bio          sql.NullString
	Status       string
	IsVerified   bool
	LastSeenAt   sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}