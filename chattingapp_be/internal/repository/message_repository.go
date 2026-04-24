package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) CreateTx(ctx context.Context, tx *sql.Tx, m *models.Message) (int64, error) {
	query := `
		INSERT INTO messages (
			conversation_id, sender_user_id, message_type, content,
			reply_to_message_id, forwarded_from_message_id, status,
			is_edited, edited_at, is_deleted, deleted_at,
			client_message_id, sent_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7,
			$8, $9, $10, $11,
			$12, $13, $14, $15
		)
		RETURNING id
	`

	var id int64
	err := tx.QueryRowContext(
		ctx,
		query,
		m.ConversationID,
		m.SenderUserID,
		m.MessageType,
		m.Content,
		m.ReplyToMessageID,
		m.ForwardedFromMessageID,
		m.Status,
		m.IsEdited,
		m.EditedAt,
		m.IsDeleted,
		m.DeletedAt,
		m.ClientMessageID,
		m.SentAt,
		m.CreatedAt,
		m.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *MessageRepository) GetByID(ctx context.Context, id int64) (*models.Message, error) {
	query := `
		SELECT id, conversation_id, sender_user_id, message_type, content,
		       reply_to_message_id, forwarded_from_message_id, status,
		       is_edited, edited_at, is_deleted, deleted_at,
		       client_message_id, sent_at, created_at, updated_at
		FROM messages
		WHERE id = $1
		LIMIT 1
	`

	var m models.Message
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID,
		&m.ConversationID,
		&m.SenderUserID,
		&m.MessageType,
		&m.Content,
		&m.ReplyToMessageID,
		&m.ForwardedFromMessageID,
		&m.Status,
		&m.IsEdited,
		&m.EditedAt,
		&m.IsDeleted,
		&m.DeletedAt,
		&m.ClientMessageID,
		&m.SentAt,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r *MessageRepository) ListByConversationID(ctx context.Context, conversationID int64, limit, offset int) ([]models.Message, error) {
	query := `
		SELECT id, conversation_id, sender_user_id, message_type, content,
		       reply_to_message_id, forwarded_from_message_id, status,
		       is_edited, edited_at, is_deleted, deleted_at,
				is_recalled, recalled_at,
				client_message_id, sent_at, created_at, updated_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY sent_at DESC, id DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(
			&m.ID,
			&m.ConversationID,
			&m.SenderUserID,
			&m.MessageType,
			&m.Content,
			&m.ReplyToMessageID,
			&m.ForwardedFromMessageID,
			&m.Status,
			&m.IsEdited,
			&m.EditedAt,
			&m.IsDeleted,
			&m.DeletedAt,
			&m.IsRecalled,
			&m.RecalledAt,
			&m.ClientMessageID,
			&m.SentAt,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, rows.Err()
}

func (r *MessageRepository) ListAscByConversationID(ctx context.Context, conversationID int64, limit, offset int) ([]models.Message, error) {
	query := `
		SELECT id, conversation_id, sender_user_id, message_type, content,
		       reply_to_message_id, forwarded_from_message_id, status,
		       is_edited, edited_at, is_deleted, deleted_at,
				is_recalled, recalled_at,
				client_message_id, sent_at, created_at, updated_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY sent_at ASC, id ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(
			&m.ID,
			&m.ConversationID,
			&m.SenderUserID,
			&m.MessageType,
			&m.Content,
			&m.ReplyToMessageID,
			&m.ForwardedFromMessageID,
			&m.Status,
			&m.IsEdited,
			&m.EditedAt,
			&m.IsDeleted,
			&m.DeletedAt,
			&m.IsRecalled,
			&m.RecalledAt,
			&m.ClientMessageID,
			&m.SentAt,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, rows.Err()
}

func (r *MessageRepository) CountByConversationID(ctx context.Context, conversationID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM messages
		WHERE conversation_id = $1
	`

	var total int
	err := r.db.QueryRowContext(ctx, query, conversationID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *MessageRepository) UpdateContent(ctx context.Context, messageID int64, content string) error {
	query := `
		UPDATE messages
		SET content = $1,
			is_edited = true,
			edited_at = $2,
			updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query, content, time.Now(), time.Now(), messageID)
	return err
}

func (r *MessageRepository) SoftDelete(ctx context.Context, messageID int64) error {
	query := `
		UPDATE messages
		SET is_deleted = true,
			deleted_at = $1,
			updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), time.Now(), messageID)
	return err
}

func (r *MessageRepository) HardDelete(ctx context.Context, messageID int64) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, messageID)
	return err
}
func (r *MessageRepository) SearchByConversationID(
	ctx context.Context,
	conversationID int64,
	keyword string,
	limit, offset int,
) ([]models.Message, error) {
	query := `
		SELECT id, conversation_id, sender_user_id, message_type, content,
		       reply_to_message_id, forwarded_from_message_id, status,
		       is_edited, edited_at, is_deleted, deleted_at,
				is_recalled, recalled_at,
				client_message_id, sent_at, created_at, updated_at
		FROM messages
		WHERE conversation_id = $1
		  AND is_deleted = false
		  AND content ILIKE '%' || $2 || '%'
		ORDER BY sent_at DESC, id DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, conversationID, keyword, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(
			&m.ID, &m.ConversationID, &m.SenderUserID, &m.MessageType, &m.Content,
			&m.ReplyToMessageID, &m.ForwardedFromMessageID, &m.Status,
			&m.IsEdited, &m.EditedAt, &m.IsDeleted, &m.DeletedAt,
			&m.IsRecalled, &m.RecalledAt,
			&m.ClientMessageID, &m.SentAt, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, rows.Err()
}

func (r *MessageRepository) ListBeforeIDByConversationID(
	ctx context.Context,
	conversationID int64,
	beforeID int64,
	limit int,
) ([]models.Message, error) {
	query := `
		SELECT id, conversation_id, sender_user_id, message_type, content,
		       reply_to_message_id, forwarded_from_message_id, status,
		       is_edited, edited_at, is_deleted, deleted_at,
				is_recalled, recalled_at,
				client_message_id, sent_at, created_at, updated_at
		FROM messages
		WHERE conversation_id = $1
		  AND id < $2
		ORDER BY id DESC
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, conversationID, beforeID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(
			&m.ID, &m.ConversationID, &m.SenderUserID, &m.MessageType, &m.Content,
			&m.ReplyToMessageID, &m.ForwardedFromMessageID, &m.Status,
			&m.IsEdited, &m.EditedAt, &m.IsDeleted, &m.DeletedAt,
			&m.IsRecalled, &m.RecalledAt,
			&m.ClientMessageID, &m.SentAt, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, rows.Err()
}

func (r *MessageRepository) UpsertReaction(
	ctx context.Context,
	messageID, userID int64,
	reactionType string,
) (*models.MessageReaction, error) {
	now := time.Now()

	updateQuery := `
		UPDATE message_reactions
		SET reaction_type = $1,
		    updated_at = $2
		WHERE message_id = $3 AND user_id = $4
		RETURNING id, message_id, user_id, reaction_type, created_at, updated_at
	`

	var reaction models.MessageReaction
	err := r.db.QueryRowContext(ctx, updateQuery, reactionType, now, messageID, userID).Scan(
		&reaction.ID,
		&reaction.MessageID,
		&reaction.UserID,
		&reaction.ReactionType,
		&reaction.CreatedAt,
		&reaction.UpdatedAt,
	)

	if err == nil {
		return &reaction, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	insertQuery := `
		INSERT INTO message_reactions (
			message_id, user_id, reaction_type, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, message_id, user_id, reaction_type, created_at, updated_at
	`

	err = r.db.QueryRowContext(ctx, insertQuery, messageID, userID, reactionType, now, now).Scan(
		&reaction.ID,
		&reaction.MessageID,
		&reaction.UserID,
		&reaction.ReactionType,
		&reaction.CreatedAt,
		&reaction.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &reaction, nil
}

func (r *MessageRepository) DeleteReaction(ctx context.Context, messageID, userID int64) error {
	query := `
		DELETE FROM message_reactions
		WHERE message_id = $1 AND user_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, messageID, userID)
	return err
}

func (r *MessageRepository) ListReactionsByMessageID(
	ctx context.Context,
	messageID int64,
) ([]models.MessageReaction, error) {
	query := `
		SELECT id, message_id, user_id, reaction_type, created_at, updated_at
		FROM message_reactions
		WHERE message_id = $1
		ORDER BY created_at ASC, id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.MessageReaction
	for rows.Next() {
		var reaction models.MessageReaction
		if err := rows.Scan(
			&reaction.ID,
			&reaction.MessageID,
			&reaction.UserID,
			&reaction.ReactionType,
			&reaction.CreatedAt,
			&reaction.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, reaction)
	}

	return result, rows.Err()
}
func (r *MessageRepository) Recall(ctx context.Context, messageID int64) error {
	query := `
		UPDATE messages
		SET is_recalled = true,
			recalled_at = $1,
			content = NULL,
			updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now, messageID)
	return err
}