package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type MessageCiphertextRepository struct {
	db *sql.DB
}

func NewMessageCiphertextRepository(db *sql.DB) *MessageCiphertextRepository {
	return &MessageCiphertextRepository{db: db}
}

func (r *MessageCiphertextRepository) Create(ctx context.Context, c *models.MessageCiphertext) (int64, error) {
	query := `
		INSERT INTO message_ciphertexts (
			message_id, target_device_id, sender_device_id,
			ciphertext, encryption_header, nonce,
			algorithm, message_version, is_delivered,
			delivered_at, created_at
		) VALUES (
			$1, $2, $3,
			$4, $5, $6,
			$7, $8, $9,
			$10, $11
		)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		c.MessageID,
		c.TargetDeviceID,
		c.SenderDeviceID,
		c.Ciphertext,
		c.EncryptionHeader,
		c.Nonce,
		c.Algorithm,
		c.MessageVersion,
		c.IsDelivered,
		c.DeliveredAt,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *MessageCiphertextRepository) ListByMessageID(ctx context.Context, messageID int64) ([]models.MessageCiphertext, error) {
	query := `
		SELECT id, message_id, target_device_id, sender_device_id,
		       ciphertext, encryption_header, nonce, algorithm,
		       message_version, is_delivered, delivered_at, created_at
		FROM message_ciphertexts
		WHERE message_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.MessageCiphertext
	for rows.Next() {
		var c models.MessageCiphertext
		if err := rows.Scan(
			&c.ID,
			&c.MessageID,
			&c.TargetDeviceID,
			&c.SenderDeviceID,
			&c.Ciphertext,
			&c.EncryptionHeader,
			&c.Nonce,
			&c.Algorithm,
			&c.MessageVersion,
			&c.IsDelivered,
			&c.DeliveredAt,
			&c.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, rows.Err()
}

func (r *MessageCiphertextRepository) ListUndeliveredByTargetDeviceID(ctx context.Context, targetDeviceID int64) ([]models.MessageCiphertext, error) {
	query := `
		SELECT id, message_id, target_device_id, sender_device_id,
		       ciphertext, encryption_header, nonce, algorithm,
		       message_version, is_delivered, delivered_at, created_at
		FROM message_ciphertexts
		WHERE target_device_id = $1
		  AND is_delivered = FALSE
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, targetDeviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.MessageCiphertext
	for rows.Next() {
		var c models.MessageCiphertext
		if err := rows.Scan(
			&c.ID,
			&c.MessageID,
			&c.TargetDeviceID,
			&c.SenderDeviceID,
			&c.Ciphertext,
			&c.EncryptionHeader,
			&c.Nonce,
			&c.Algorithm,
			&c.MessageVersion,
			&c.IsDelivered,
			&c.DeliveredAt,
			&c.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, rows.Err()
}

func (r *MessageCiphertextRepository) MarkDelivered(ctx context.Context, id int64) error {
	query := `
		UPDATE message_ciphertexts
		SET is_delivered = TRUE,
			delivered_at = $2
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id, time.Now())
	return err
}
func (r *MessageCiphertextRepository) CreateTx(ctx context.Context, tx *sql.Tx, c *models.MessageCiphertext) (int64, error) {
	query := `
		INSERT INTO message_ciphertexts (
			message_id, target_device_id, sender_device_id,
			ciphertext, encryption_header, nonce,
			algorithm, message_version, is_delivered,
			delivered_at, created_at
		) VALUES (
			$1, $2, $3,
			$4, $5, $6,
			$7, $8, $9,
			$10, $11
		)
		RETURNING id
	`

	var id int64
	err := tx.QueryRowContext(
		ctx,
		query,
		c.MessageID,
		c.TargetDeviceID,
		c.SenderDeviceID,
		c.Ciphertext,
		c.EncryptionHeader,
		c.Nonce,
		c.Algorithm,
		c.MessageVersion,
		c.IsDelivered,
		c.DeliveredAt,
		c.CreatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *MessageCiphertextRepository) GetByID(ctx context.Context, id int64) (*models.MessageCiphertext, error) {
	query := `
		SELECT id, message_id, target_device_id, sender_device_id,
		       ciphertext, encryption_header, nonce, algorithm,
		       message_version, is_delivered, delivered_at, created_at
		FROM message_ciphertexts
		WHERE id = $1
		LIMIT 1
	`

	var c models.MessageCiphertext
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.MessageID,
		&c.TargetDeviceID,
		&c.SenderDeviceID,
		&c.Ciphertext,
		&c.EncryptionHeader,
		&c.Nonce,
		&c.Algorithm,
		&c.MessageVersion,
		&c.IsDelivered,
		&c.DeliveredAt,
		&c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}