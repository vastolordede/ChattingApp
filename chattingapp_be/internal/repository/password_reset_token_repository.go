package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
)

type PasswordResetTokenRepository struct {
	db *sql.DB
}

func NewPasswordResetTokenRepository(db *sql.DB) *PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{db: db}
}

func (r *PasswordResetTokenRepository) Create(ctx context.Context, t *models.PasswordResetToken) (int64, error) {
	query := `
		INSERT INTO password_reset_tokens (
			user_id, token_hash, expires_at, used_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		t.UserID,
		t.TokenHash,
		t.ExpiresAt,
		t.UsedAt,
		t.CreatedAt,
		t.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PasswordResetTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, used_at, created_at, updated_at
		FROM password_reset_tokens
		WHERE token_hash = $1
		LIMIT 1
	`

	var t models.PasswordResetToken
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&t.ID,
		&t.UserID,
		&t.TokenHash,
		&t.ExpiresAt,
		&t.UsedAt,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *PasswordResetTokenRepository) MarkUsed(ctx context.Context, id int64, usedAt sql.NullTime) error {
	query := `
		UPDATE password_reset_tokens
		SET used_at = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, usedAt, id)
	return err
}