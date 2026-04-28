package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type ConversationMemberRepository struct {
	db *sql.DB
}

func NewConversationMemberRepository(db *sql.DB) *ConversationMemberRepository {
	return &ConversationMemberRepository{db: db}
}

func (r *ConversationMemberRepository) CreateTx(ctx context.Context, tx *sql.Tx, cm *models.ConversationMember) (int64, error) {
	query := `
		INSERT INTO conversation_members (
			conversation_id, user_id, role, joined_at, left_at, is_active, nickname,
			last_read_message_id, last_read_at, mute_until, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12
		)
		RETURNING id
	`

	var id int64
	err := tx.QueryRowContext(
		ctx,
		query,
		cm.ConversationID,
		cm.UserID,
		cm.Role,
		cm.JoinedAt,
		cm.LeftAt,
		cm.IsActive,
		cm.Nickname,
		cm.LastReadMessageID,
		cm.LastReadAt,
		cm.MuteUntil,
		cm.CreatedAt,
		cm.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ConversationMemberRepository) GetByConversationAndUser(ctx context.Context, conversationID, userID int64) (*models.ConversationMember, error) {
	query := `
		SELECT id, conversation_id, user_id, role, joined_at, left_at, is_active, nickname,
		       last_read_message_id, last_read_at, mute_until, created_at, updated_at
		FROM conversation_members
		WHERE conversation_id = $1 AND user_id = $2
		LIMIT 1
	`

	var cm models.ConversationMember
	err := r.db.QueryRowContext(ctx, query, conversationID, userID).Scan(
		&cm.ID,
		&cm.ConversationID,
		&cm.UserID,
		&cm.Role,
		&cm.JoinedAt,
		&cm.LeftAt,
		&cm.IsActive,
		&cm.Nickname,
		&cm.LastReadMessageID,
		&cm.LastReadAt,
		&cm.MuteUntil,
		&cm.CreatedAt,
		&cm.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cm, nil
}

func (r *ConversationMemberRepository) ListByConversationID(ctx context.Context, conversationID int64) ([]models.ConversationMember, error) {
	query := `
		SELECT id, conversation_id, user_id, role, joined_at, left_at, is_active, nickname,
		       last_read_message_id, last_read_at, mute_until, created_at, updated_at
		FROM conversation_members
		WHERE conversation_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.ConversationMember
	for rows.Next() {
		var cm models.ConversationMember
		if err := rows.Scan(
			&cm.ID,
			&cm.ConversationID,
			&cm.UserID,
			&cm.Role,
			&cm.JoinedAt,
			&cm.LeftAt,
			&cm.IsActive,
			&cm.Nickname,
			&cm.LastReadMessageID,
			&cm.LastReadAt,
			&cm.MuteUntil,
			&cm.CreatedAt,
			&cm.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, cm)
	}

	return result, rows.Err()
}

func (r *ConversationMemberRepository) UpdateLastRead(ctx context.Context, conversationID, userID, messageID int64) error {
	query := `
		UPDATE conversation_members
		SET last_read_message_id = $1,
			last_read_at = $2,
			updated_at = $3
		WHERE conversation_id = $4
		  AND user_id = $5
	`

	_, err := r.db.ExecContext(ctx, query, messageID, time.Now(), time.Now(), conversationID, userID)
	return err
}

func (r *ConversationMemberRepository) UpdateNickname(ctx context.Context, conversationID, userID int64, nickname sql.NullString) error {
	query := `
		UPDATE conversation_members
		SET nickname = $1,
			updated_at = $2
		WHERE conversation_id = $3
		  AND user_id = $4
	`

	_, err := r.db.ExecContext(ctx, query, nickname, time.Now(), conversationID, userID)
	return err
}

func (r *ConversationMemberRepository) UpdateMuteUntil(ctx context.Context, conversationID, userID int64, muteUntil sql.NullTime) error {
	query := `
		UPDATE conversation_members
		SET mute_until = $1,
			updated_at = $2
		WHERE conversation_id = $3
		  AND user_id = $4
	`

	_, err := r.db.ExecContext(ctx, query, muteUntil, time.Now(), conversationID, userID)
	return err
}
func (r *ConversationMemberRepository) UpdatePinned(
	ctx context.Context,
	conversationID, userID int64,
	isPinned bool,
) error {
	query := `
		UPDATE conversation_members
		SET is_pinned = $1,
			updated_at = $2
		WHERE conversation_id = $3 AND user_id = $4
	`
	_, err := r.db.ExecContext(ctx, query, isPinned, time.Now(), conversationID, userID)
	return err
}

func (r *ConversationMemberRepository) UpdateArchived(
	ctx context.Context,
	conversationID, userID int64,
	isArchived bool,
) error {
	query := `
		UPDATE conversation_members
		SET is_archived = $1,
			updated_at = $2
		WHERE conversation_id = $3 AND user_id = $4
	`
	_, err := r.db.ExecContext(ctx, query, isArchived, time.Now(), conversationID, userID)
	return err
}
func (r *ConversationMemberRepository) CountUnread(
	ctx context.Context,
	userID int64,
) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM conversation_members cm
		JOIN conversations c ON c.id = cm.conversation_id
		WHERE cm.user_id = $1
		  AND c.last_message_id IS NOT NULL
		  AND (cm.last_read_message_id IS NULL OR cm.last_read_message_id < c.last_message_id)
	`
	var count int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}
