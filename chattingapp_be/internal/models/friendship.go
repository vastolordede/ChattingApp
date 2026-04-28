package models

import (
	"database/sql"
	"time"
)

type Friendship struct {
	ID        int64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	EndedAt   sql.NullTime
}
