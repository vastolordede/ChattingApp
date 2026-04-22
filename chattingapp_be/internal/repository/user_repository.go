package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) (int64, error) {
	query := `
		INSERT INTO users (
			username, full_name, email, phone_number, password_hash,
			avatar_url, bio, status, is_verified, last_seen_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11, $12
		)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		u.Username,
		u.FullName,
		u.Email,
		u.PhoneNumber,
		u.PasswordHash,
		u.AvatarURL,
		u.Bio,
		u.Status,
		u.IsVerified,
		u.LastSeenAt,
		u.CreatedAt,
		u.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, username, full_name, email, phone_number, password_hash,
		       avatar_url, bio, status, is_verified, last_seen_at, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	var u models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Username,
		&u.FullName,
		&u.Email,
		&u.PhoneNumber,
		&u.PasswordHash,
		&u.AvatarURL,
		&u.Bio,
		&u.Status,
		&u.IsVerified,
		&u.LastSeenAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, full_name, email, phone_number, password_hash,
		       avatar_url, bio, status, is_verified, last_seen_at, created_at, updated_at
		FROM users
		WHERE username = $1
		LIMIT 1
	`

	var u models.User
	err := r.db.QueryRowContext(ctx, query, strings.TrimSpace(username)).Scan(
		&u.ID,
		&u.Username,
		&u.FullName,
		&u.Email,
		&u.PhoneNumber,
		&u.PasswordHash,
		&u.AvatarURL,
		&u.Bio,
		&u.Status,
		&u.IsVerified,
		&u.LastSeenAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, full_name, email, phone_number, password_hash,
		       avatar_url, bio, status, is_verified, last_seen_at, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
		LIMIT 1
	`

	var u models.User
	err := r.db.QueryRowContext(ctx, query, strings.TrimSpace(email)).Scan(
		&u.ID,
		&u.Username,
		&u.FullName,
		&u.Email,
		&u.PhoneNumber,
		&u.PasswordHash,
		&u.AvatarURL,
		&u.Bio,
		&u.Status,
		&u.IsVerified,
		&u.LastSeenAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	query := `
		SELECT id, username, full_name, email, phone_number, password_hash,
		       avatar_url, bio, status, is_verified, last_seen_at, created_at, updated_at
		FROM users
		WHERE phone_number = $1
		LIMIT 1
	`

	var u models.User
	err := r.db.QueryRowContext(ctx, query, strings.TrimSpace(phoneNumber)).Scan(
		&u.ID,
		&u.Username,
		&u.FullName,
		&u.Email,
		&u.PhoneNumber,
		&u.PasswordHash,
		&u.AvatarURL,
		&u.Bio,
		&u.Status,
		&u.IsVerified,
		&u.LastSeenAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetByIdentifier(ctx context.Context, identifier string) (*models.User, error) {
	query := `
		SELECT id, username, full_name, email, phone_number, password_hash,
		       avatar_url, bio, status, is_verified, last_seen_at, created_at, updated_at
		FROM users
		WHERE username = $1
		   OR LOWER(email) = LOWER($1)
		   OR phone_number = $1
		LIMIT 1
	`

	identifier = strings.TrimSpace(identifier)

	var u models.User
	err := r.db.QueryRowContext(ctx, query, identifier).Scan(
		&u.ID,
		&u.Username,
		&u.FullName,
		&u.Email,
		&u.PhoneNumber,
		&u.PasswordHash,
		&u.AvatarURL,
		&u.Bio,
		&u.Status,
		&u.IsVerified,
		&u.LastSeenAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) UpdateProfile(
	ctx context.Context,
	userID int64,
	fullName string,
	avatarURL sql.NullString,
	bio sql.NullString,
) error {
	query := `
		UPDATE users
		SET full_name = $1,
			avatar_url = $2,
			bio = $3,
			updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		strings.TrimSpace(fullName),
		avatarURL,
		bio,
		time.Now(),
		userID,
	)

	return err
}

func (r *UserRepository) SearchUsers(ctx context.Context, keyword string, limit, offset int) ([]models.User, error) {
	query := `
		SELECT id, username, full_name, email, phone_number, password_hash,
		       avatar_url, bio, status, is_verified, last_seen_at, created_at, updated_at
		FROM users
		WHERE username ILIKE '%' || $1 || '%'
		   OR full_name ILIKE '%' || $1 || '%'
		   OR email ILIKE '%' || $1 || '%'
		   OR phone_number ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, strings.TrimSpace(keyword), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.FullName,
			&u.Email,
			&u.PhoneNumber,
			&u.PasswordHash,
			&u.AvatarURL,
			&u.Bio,
			&u.Status,
			&u.IsVerified,
			&u.LastSeenAt,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, u)
	}

	return result, rows.Err()
}

func (r *UserRepository) CountSearchUsers(ctx context.Context, keyword string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM users
		WHERE username ILIKE '%' || $1 || '%'
		   OR full_name ILIKE '%' || $1 || '%'
		   OR email ILIKE '%' || $1 || '%'
		   OR phone_number ILIKE '%' || $1 || '%'
	`

	var total int
	err := r.db.QueryRowContext(ctx, query, strings.TrimSpace(keyword)).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
func (r *UserRepository) UpdatePasswordHash(ctx context.Context, userID int64, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1,
			updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, passwordHash, time.Now(), userID)
	return err
}