package dto

type AddAttachmentsToMessageRequest struct {
	MessageID    int64                            `json:"message_id" validate:"required"`
	Attachments  []CreateMessageAttachmentRequest `json:"attachments" validate:"required"`
}