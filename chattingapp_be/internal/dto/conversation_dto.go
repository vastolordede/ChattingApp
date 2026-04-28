package dto

type CreateDirectConversationRequest struct {
	TargetUserID int64 `json:"target_user_id" validate:"required"`
}

type ConversationMemberResponse struct {
	UserID      int64   `json:"user_id"`
	Role        string  `json:"role"`
	Nickname    *string `json:"nickname,omitempty"`
	IsActive    bool    `json:"is_active"`
	JoinedAt    string  `json:"joined_at"`
	LeftAt      *string `json:"left_at,omitempty"`
	LastReadAt  *string `json:"last_read_at,omitempty"`
	LastReadMsg *int64  `json:"last_read_message_id,omitempty"`
}

type ConversationListItemResponse struct {
	ID               int64                   `json:"id"`
	ConversationType string                  `json:"conversation_type"`
	Title            *string                 `json:"title,omitempty"`
	AvatarURL        *string                 `json:"avatar_url,omitempty"`
	Status           string                  `json:"status"`
	LastMessageID    *int64                  `json:"last_message_id,omitempty"`
	LastMessageAt    *string                 `json:"last_message_at,omitempty"`
	LastMessage      *MessagePreviewResponse `json:"last_message,omitempty"`
	CreatedAt        string                  `json:"created_at"`
	UpdatedAt        string                  `json:"updated_at"`
}

type ConversationDetailResponse struct {
	ID               int64                       `json:"id"`
	ConversationType string                      `json:"conversation_type"`
	Title            *string                     `json:"title,omitempty"`
	AvatarURL        *string                     `json:"avatar_url,omitempty"`
	Status           string                      `json:"status"`
	CreatedByUserID  *int64                      `json:"created_by_user_id,omitempty"`
	LastMessageID    *int64                      `json:"last_message_id,omitempty"`
	LastMessageAt    *string                     `json:"last_message_at,omitempty"`
	CreatedAt        string                      `json:"created_at"`
	UpdatedAt        string                      `json:"updated_at"`
	Members          []ConversationMemberSummary `json:"members"`
}

type ConversationMemberSummary struct {
	UserID      int64   `json:"user_id"`
	Username    string  `json:"username"`
	FullName    string  `json:"full_name"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Role        string  `json:"role"`
	Nickname    *string `json:"nickname,omitempty"`
	IsActive    bool    `json:"is_active"`
	JoinedAt    string  `json:"joined_at"`
	LastReadAt  *string `json:"last_read_at,omitempty"`
	LastReadMsg *int64  `json:"last_read_message_id,omitempty"`
}

type MarkConversationReadRequest struct {
	LastReadMessageID int64 `json:"last_read_message_id" validate:"required"`
}
type TypingRequest struct {
	IsTyping bool `json:"is_typing"`
}
