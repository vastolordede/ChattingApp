package dto

type FriendListItemResponse struct {
	FriendshipID int64             `json:"friendship_id"`
	Status       string            `json:"status"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
	EndedAt      *string           `json:"ended_at,omitempty"`
	Friend       FriendUserSummary `json:"friend"`
}

type RemoveFriendRequest struct {
	FriendUserID int64 `json:"friend_user_id" validate:"required"`
}

type BlockFriendRequest struct {
	FriendUserID int64 `json:"friend_user_id" validate:"required"`
}
