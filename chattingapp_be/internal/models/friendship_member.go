package models

import "time"

type FriendshipMember struct {
	ID           int64
	FriendshipID int64
	UserID       int64
	CreatedAt    time.Time
}
