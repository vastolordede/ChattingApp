package dto

type UpdateConversationNicknameRequest struct {
	Nickname *string `json:"nickname,omitempty"`
}

type MuteConversationRequest struct {
	MuteUntil *string `json:"mute_until,omitempty"`
}