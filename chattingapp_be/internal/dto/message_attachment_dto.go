package dto

type CreateMessageAttachmentRequest struct {
	MessageID         int64   `json:"message_id" validate:"required"`
	AttachmentType    string  `json:"attachment_type" validate:"required"` // image/video/file/audio/link
	FileName          string  `json:"file_name" validate:"required"`
	MimeType          string  `json:"mime_type" validate:"required"`
	FileSize          int64   `json:"file_size" validate:"required"`
	FileURL           string  `json:"file_url" validate:"required"`
	ThumbnailURL      *string `json:"thumbnail_url,omitempty"`
	Width             *int    `json:"width,omitempty"`
	Height            *int    `json:"height,omitempty"`
	DurationSeconds   *int    `json:"duration_seconds,omitempty"`
	Checksum          *string `json:"checksum,omitempty"`
	EncryptionKeyHint *string `json:"encryption_key_hint,omitempty"`
}

type MessageAttachmentResponse struct {
	ID                int64   `json:"id"`
	AttachmentType    string  `json:"attachment_type"`
	FileName          string  `json:"file_name"`
	MimeType          string  `json:"mime_type"`
	FileSize          int64   `json:"file_size"`
	FileURL           string  `json:"file_url"`
	ThumbnailURL      *string `json:"thumbnail_url,omitempty"`
	Width             *int    `json:"width,omitempty"`
	Height            *int    `json:"height,omitempty"`
	DurationSeconds   *int    `json:"duration_seconds,omitempty"`
	Checksum          *string `json:"checksum,omitempty"`
	EncryptionKeyHint *string `json:"encryption_key_hint,omitempty"`
	CreatedAt         string  `json:"created_at"`
}
