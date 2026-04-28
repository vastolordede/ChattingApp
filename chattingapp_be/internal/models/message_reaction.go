package models

import "time"

type MessageReaction struct {
	ID           int64
	MessageID    int64
	UserID       int64
	ReactionType string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
