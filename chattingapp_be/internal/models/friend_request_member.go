package models

import "time"

type FriendRequestMember struct {
	ID              int64
	FriendRequestID int64
	UserID          int64
	Role            string
	CreatedAt       time.Time
}