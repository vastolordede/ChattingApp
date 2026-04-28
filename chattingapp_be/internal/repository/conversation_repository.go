package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type ConversationRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) CreateTx(ctx context.Context, tx *sql.Tx, c *models.Conversation) (int64, error) {
	query := `
		INSERT INTO conversations (
			conversation_type, title, avatar_url, created_by_user_id, status,
			last_message_id, last_message_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	var id int64
	err := tx.QueryRowContext(
		ctx,
		query,
		c.ConversationType,
		c.Title,
		c.AvatarURL,
		c.CreatedByUserID,
		c.Status,
		c.LastMessageID,
		c.LastMessageAt,
		c.CreatedAt,
		c.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ConversationRepository) GetByID(ctx context.Context, id int64) (*models.Conversation, error) {
	query := `
		SELECT id, conversation_type, title, avatar_url, created_by_user_id, status,
		       last_message_id, last_message_at, created_at, updated_at
		FROM conversations
		WHERE id = $1
		LIMIT 1
	`

	var c models.Conversation
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.ConversationType,
		&c.Title,
		&c.AvatarURL,
		&c.CreatedByUserID,
		&c.Status,
		&c.LastMessageID,
		&c.LastMessageAt,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

func (r *ConversationRepository) FindDirectConversationBetweenUsers(ctx context.Context, userAID, userBID int64) (*models.Conversation, error) {
	query := `
		SELECT c.id, c.conversation_type, c.title, c.avatar_url, c.created_by_user_id, c.status,
		       c.last_message_id, c.last_message_at, c.created_at, c.updated_at
		FROM conversations c
		JOIN conversation_members cm1 ON cm1.conversation_id = c.id
		JOIN conversation_members cm2 ON cm2.conversation_id = c.id
		WHERE c.conversation_type = 'direct'
		  AND c.status = 'active'
		  AND cm1.user_id = $1
		  AND cm2.user_id = $2
		  AND cm1.is_active = true
		  AND cm2.is_active = true
		LIMIT 1
	`

	var c models.Conversation
	err := r.db.QueryRowContext(ctx, query, userAID, userBID).Scan(
		&c.ID,
		&c.ConversationType,
		&c.Title,
		&c.AvatarURL,
		&c.CreatedByUserID,
		&c.Status,
		&c.LastMessageID,
		&c.LastMessageAt,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

func (r *ConversationRepository) UpdateLastMessageTx(
	ctx context.Context,
	tx *sql.Tx,
	conversationID int64,
	lastMessageID int64,
	lastMessageAt time.Time,
) error {
	query := `
		UPDATE conversations
		SET last_message_id = $1,
			last_message_at = $2,
			updated_at = $3
		WHERE id = $4
	`

	_, err := tx.ExecContext(ctx, query, lastMessageID, lastMessageAt, time.Now(), conversationID)
	return err
}

func (r *ConversationRepository) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]models.Conversation, error) {
	query := `
		SELECT c.id, c.conversation_type, c.title, c.avatar_url, c.created_by_user_id, c.status,
		       c.last_message_id, c.last_message_at, c.created_at, c.updated_at
		FROM conversations c
		JOIN conversation_members cm ON cm.conversation_id = c.id
		WHERE cm.user_id = $1
		  AND cm.is_active = true
		ORDER BY c.last_message_at DESC NULLS LAST, c.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Conversation
	for rows.Next() {
		var c models.Conversation
		if err := rows.Scan(
			&c.ID,
			&c.ConversationType,
			&c.Title,
			&c.AvatarURL,
			&c.CreatedByUserID,
			&c.Status,
			&c.LastMessageID,
			&c.LastMessageAt,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, rows.Err()
}

func (r *ConversationRepository) CountByUserID(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM conversations c
		JOIN conversation_members cm ON cm.conversation_id = c.id
		WHERE cm.user_id = $1
		  AND cm.is_active = true
	`

	var total int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
