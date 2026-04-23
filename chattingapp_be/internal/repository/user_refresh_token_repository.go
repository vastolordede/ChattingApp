package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
	"fmt"
)

type UserRefreshTokenRepository struct {
	db *sql.DB
}

func NewUserRefreshTokenRepository(db *sql.DB) *UserRefreshTokenRepository {
	return &UserRefreshTokenRepository{db: db}
}

func (r *UserRefreshTokenRepository) Create(ctx context.Context, token *models.UserRefreshToken) (int64, error) {
	query := `
		INSERT INTO user_refresh_tokens (
			user_id,
			user_device_id,
			token_hash,
			expires_at,
			revoked_at,
			last_used_at,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
fmt.Println("REPO: Create refresh token called")
	fmt.Println("REPO: user_id =", token.UserID)
	fmt.Println("REPO: user_device_id =", token.UserDeviceID)
	fmt.Println("REPO: token_hash =", token.TokenHash)
	fmt.Println("REPO: expires_at =", token.ExpiresAt)
	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		token.UserID,
		token.UserDeviceID,
		token.TokenHash,
		token.ExpiresAt,
		token.RevokedAt,
		token.LastUsedAt,
		token.CreatedAt,
		token.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.UserRefreshToken, error) {
	query := `
		SELECT
			id,
			user_id,
			user_device_id,
			token_hash,
			expires_at,
			revoked_at,
			last_used_at,
			created_at,
			updated_at
		FROM user_refresh_tokens
		WHERE token_hash = $1
		LIMIT 1
	`

	var t models.UserRefreshToken
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&t.ID,
		&t.UserID,
		&t.UserDeviceID,
		&t.TokenHash,
		&t.ExpiresAt,
		&t.RevokedAt,
		&t.LastUsedAt,
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

func (r *UserRefreshTokenRepository) RevokeByID(ctx context.Context, id int64, revokedAt sql.NullTime) error {
	query := `
		UPDATE user_refresh_tokens
		SET revoked_at = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, revokedAt, id)
	return err
}
func (r *UserRefreshTokenRepository) UpdateLastUsedAt(ctx context.Context, id int64, lastUsedAt sql.NullTime) error {
	query := `
		UPDATE user_refresh_tokens
		SET last_used_at = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, lastUsedAt, id)
	return err
}
func (r *UserRefreshTokenRepository) RevokeByDeviceID(ctx context.Context, deviceID int64) error {
	query := `
		UPDATE user_refresh_tokens
		SET revoked_at = NOW(), updated_at = NOW()
		WHERE user_device_id = $1 AND revoked_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, deviceID)
	return err
}