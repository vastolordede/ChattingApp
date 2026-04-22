package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
)

type MessageAttachmentRepository struct {
	db *sql.DB
}

func NewMessageAttachmentRepository(db *sql.DB) *MessageAttachmentRepository {
	return &MessageAttachmentRepository{db: db}
}

func (r *MessageAttachmentRepository) CreateTx(ctx context.Context, tx *sql.Tx, a *models.MessageAttachment) (int64, error) {
	query := `
		INSERT INTO message_attachments (
			message_id, attachment_type, file_name, mime_type, file_size, file_url,
			thumbnail_url, width, height, duration_seconds, checksum, encryption_key_hint, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13
		)
		RETURNING id
	`

	var id int64
	err := tx.QueryRowContext(
		ctx,
		query,
		a.MessageID,
		a.AttachmentType,
		a.FileName,
		a.MimeType,
		a.FileSize,
		a.FileURL,
		a.ThumbnailURL,
		a.Width,
		a.Height,
		a.DurationSeconds,
		a.Checksum,
		a.EncryptionKeyHint,
		a.CreatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *MessageAttachmentRepository) ListByMessageID(ctx context.Context, messageID int64) ([]models.MessageAttachment, error) {
	query := `
		SELECT id, message_id, attachment_type, file_name, mime_type, file_size, file_url,
		       thumbnail_url, width, height, duration_seconds, checksum, encryption_key_hint, created_at
		FROM message_attachments
		WHERE message_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.MessageAttachment
	for rows.Next() {
		var a models.MessageAttachment
		if err := rows.Scan(
			&a.ID,
			&a.MessageID,
			&a.AttachmentType,
			&a.FileName,
			&a.MimeType,
			&a.FileSize,
			&a.FileURL,
			&a.ThumbnailURL,
			&a.Width,
			&a.Height,
			&a.DurationSeconds,
			&a.Checksum,
			&a.EncryptionKeyHint,
			&a.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, rows.Err()
}