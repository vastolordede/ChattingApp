package dto

type SendFriendRequestRequest struct {
	ReceiverUserID int64   `json:"receiver_user_id" validate:"required"`
	Message        *string `json:"message,omitempty"`
}

type RespondFriendRequestRequest struct {
	Action string `json:"action" validate:"required"` // accepted | rejected | cancelled
}

type FriendRequestResponse struct {
	ID          int64             `json:"id"`
	Status      string            `json:"status"`
	Message     *string           `json:"message,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	RespondedAt *string           `json:"responded_at,omitempty"`
	ExpiredAt   *string           `json:"expired_at,omitempty"`
	Sender      FriendUserSummary `json:"sender"`
	Receiver    FriendUserSummary `json:"receiver"`
}

type FriendUserSummary struct {
	ID         int64   `json:"id"`
	Username   string  `json:"username"`
	FullName   string  `json:"full_name"`
	AvatarURL  *string `json:"avatar_url,omitempty"`
	IsVerified bool    `json:"is_verified"`
}
