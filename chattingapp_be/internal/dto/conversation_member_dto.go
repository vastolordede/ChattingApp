package dto

type UpdateConversationNicknameRequest struct {
	Nickname *string `json:"nickname,omitempty"`
}

type MuteConversationRequest struct {
	MuteUntil *string `json:"mute_until,omitempty"`
}
type PinConversationRequest struct {
	IsPinned bool `json:"is_pinned"`
}

type ArchiveConversationRequest struct {
	IsArchived bool `json:"is_archived"`
}
type UnreadCountResponse struct {
	UnreadCount int64 `json:"unread_count"`
}
