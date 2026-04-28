package dto

type SendMessageRequest struct {
	ConversationID   int64   `json:"conversation_id" validate:"required"`
	MessageType      string  `json:"message_type" validate:"required"` // text/image/video/file/audio/link/system
	Content          *string `json:"content,omitempty"`
	ReplyToMessageID *int64  `json:"reply_to_message_id,omitempty"`
	ClientMessageID  *string `json:"client_message_id,omitempty"`

	Attachments []CreateMessageAttachmentRequest `json:"attachments,omitempty"`
}

type EditMessageRequest struct {
	Content string `json:"content" validate:"required"`
}

type DeleteMessageRequest struct {
	SoftDelete bool `json:"soft_delete"`
}

type ListMessagesRequest struct {
	ConversationID int64 `json:"conversation_id"`
	Page           int   `json:"page"`
	Limit          int   `json:"limit"`
}

type MessageSenderInfo struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	FullName  string  `json:"full_name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type MessagePreviewResponse struct {
	ID           int64   `json:"id"`
	SenderUserID int64   `json:"sender_user_id"`
	MessageType  string  `json:"message_type"`
	Content      *string `json:"content,omitempty"`
	Status       string  `json:"status"`
	IsEdited     bool    `json:"is_edited"`
	IsDeleted    bool    `json:"is_deleted"`
	SentAt       string  `json:"sent_at"`
}

type MessageResponse struct {
	ID                     int64                       `json:"id"`
	ConversationID         int64                       `json:"conversation_id"`
	SenderUserID           int64                       `json:"sender_user_id"`
	MessageType            string                      `json:"message_type"`
	Content                *string                     `json:"content,omitempty"`
	ReplyToMessageID       *int64                      `json:"reply_to_message_id,omitempty"`
	ForwardedFromMessageID *int64                      `json:"forwarded_from_message_id,omitempty"`
	Status                 string                      `json:"status"`
	IsEdited               bool                        `json:"is_edited"`
	EditedAt               *string                     `json:"edited_at,omitempty"`
	IsDeleted              bool                        `json:"is_deleted"`
	DeletedAt              *string                     `json:"deleted_at,omitempty"`
	IsRecalled             bool                        `json:"is_recalled"`
	RecalledAt             *string                     `json:"recalled_at,omitempty"`
	ClientMessageID        *string                     `json:"client_message_id,omitempty"`
	SentAt                 string                      `json:"sent_at"`
	CreatedAt              string                      `json:"created_at"`
	UpdatedAt              string                      `json:"updated_at"`
	Sender                 *MessageSenderInfo          `json:"sender,omitempty"`
	ReplyToMessagePreview  *MessagePreviewResponse     `json:"reply_to_message_preview,omitempty"`
	Attachments            []MessageAttachmentResponse `json:"attachments,omitempty"`
}
type ForwardMessageRequest struct {
	TargetConversationID int64   `json:"target_conversation_id" validate:"required"`
	Content              *string `json:"content,omitempty"`
	ClientMessageID      *string `json:"client_message_id,omitempty"`
}

type ReactMessageRequest struct {
	ReactionType string `json:"reaction_type" validate:"required"` // like, love, haha, wow, sad, angry...
}

type MessageReactionResponse struct {
	ID           int64  `json:"id"`
	MessageID    int64  `json:"message_id"`
	UserID       int64  `json:"user_id"`
	ReactionType string `json:"reaction_type"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
