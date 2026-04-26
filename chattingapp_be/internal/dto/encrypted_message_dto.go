package dto

type EncryptedCiphertextItemRequest struct {
	TargetDeviceID   int64   `json:"target_device_id" validate:"required"`
	Ciphertext       string  `json:"ciphertext" validate:"required"`
	EncryptionHeader *string `json:"encryption_header,omitempty"`
	Nonce            *string `json:"nonce,omitempty"`
	Algorithm        string  `json:"algorithm,omitempty"`
	MessageVersion   int     `json:"message_version,omitempty"`
}

type SendEncryptedMessageRequest struct {
	ConversationID     int64                            `json:"conversation_id" validate:"required"`
	SenderDeviceUUID   string                           `json:"sender_device_uuid" validate:"required"`
	ReplyToMessageID   *int64                           `json:"reply_to_message_id,omitempty"`
	ClientMessageID    *string                          `json:"client_message_id,omitempty"`
	Ciphertexts        []EncryptedCiphertextItemRequest `json:"ciphertexts" validate:"required"`
}

type EncryptedCiphertextResponse struct {
	ID               int64   `json:"id"`
	MessageID        int64   `json:"message_id"`
	TargetDeviceID   int64   `json:"target_device_id"`
	SenderDeviceID   *int64  `json:"sender_device_id,omitempty"`
	Ciphertext       string  `json:"ciphertext"`
	EncryptionHeader *string `json:"encryption_header,omitempty"`
	Nonce            *string `json:"nonce,omitempty"`
	Algorithm        string  `json:"algorithm"`
	MessageVersion   int     `json:"message_version"`
	IsDelivered      bool    `json:"is_delivered"`
	DeliveredAt      *string `json:"delivered_at,omitempty"`
	CreatedAt        string  `json:"created_at"`
}

type SendEncryptedMessageResponse struct {
	Message     MessageResponse                `json:"message"`
	Ciphertexts []EncryptedCiphertextResponse `json:"ciphertexts"`
}