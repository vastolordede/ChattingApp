package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type FriendshipRepository struct {
	db *sql.DB
}

func NewFriendshipRepository(db *sql.DB) *FriendshipRepository {
	return &FriendshipRepository{db: db}
}

func (r *FriendshipRepository) CreateTx(ctx context.Context, tx *sql.Tx, fs *models.Friendship) (int64, error) {
	query := `
		INSERT INTO friendships (status, created_at, updated_at, ended_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int64
	err := tx.QueryRowContext(
		ctx,
		query,
		fs.Status,
		fs.CreatedAt,
		fs.UpdatedAt,
		fs.EndedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *FriendshipRepository) AddMemberTx(ctx context.Context, tx *sql.Tx, member *models.FriendshipMember) error {
	query := `
		INSERT INTO friendship_members (friendship_id, user_id, created_at)
		VALUES ($1, $2, $3)
	`

	_, err := tx.ExecContext(ctx, query, member.FriendshipID, member.UserID, member.CreatedAt)
	return err
}

func (r *FriendshipRepository) FindActiveByUserPair(ctx context.Context, userAID, userBID int64) (*models.Friendship, error) {
	query := `
		SELECT f.id, f.status, f.created_at, f.updated_at, f.ended_at
		FROM friendships f
		JOIN friendship_members fm1 ON fm1.friendship_id = f.id
		JOIN friendship_members fm2 ON fm2.friendship_id = f.id
		WHERE f.status = 'active'
		  AND fm1.user_id = $1
		  AND fm2.user_id = $2
		  AND fm1.user_id <> fm2.user_id
		LIMIT 1
	`

	var fs models.Friendship
	err := r.db.QueryRowContext(ctx, query, userAID, userBID).Scan(
		&fs.ID,
		&fs.Status,
		&fs.CreatedAt,
		&fs.UpdatedAt,
		&fs.EndedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &fs, nil
}

func (r *FriendshipRepository) GetByID(ctx context.Context, friendshipID int64) (*models.Friendship, error) {
	query := `
		SELECT id, status, created_at, updated_at, ended_at
		FROM friendships
		WHERE id = $1
		LIMIT 1
	`

	var fs models.Friendship
	err := r.db.QueryRowContext(ctx, query, friendshipID).Scan(
		&fs.ID,
		&fs.Status,
		&fs.CreatedAt,
		&fs.UpdatedAt,
		&fs.EndedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &fs, nil
}

func (r *FriendshipRepository) GetMembersByFriendshipID(ctx context.Context, friendshipID int64) ([]models.FriendshipMember, error) {
	query := `
		SELECT id, friendship_id, user_id, created_at
		FROM friendship_members
		WHERE friendship_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, friendshipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.FriendshipMember
	for rows.Next() {
		var item models.FriendshipMember
		if err := rows.Scan(
			&item.ID,
			&item.FriendshipID,
			&item.UserID,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func (r *FriendshipRepository) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]models.Friendship, error) {
	query := `
		SELECT DISTINCT f.id, f.status, f.created_at, f.updated_at, f.ended_at
		FROM friendships f
		JOIN friendship_members fm ON fm.friendship_id = f.id
		WHERE fm.user_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Friendship
	for rows.Next() {
		var item models.Friendship
		if err := rows.Scan(
			&item.ID,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.EndedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func (r *FriendshipRepository) ListFriendUserIDs(ctx context.Context, userID int64) ([]int64, error) {
	query := `
		SELECT fm2.user_id
		FROM friendships f
		JOIN friendship_members fm1 ON fm1.friendship_id = f.id
		JOIN friendship_members fm2 ON fm2.friendship_id = f.id
		WHERE f.status = 'active'
		  AND fm1.user_id = $1
		  AND fm2.user_id <> $1
		ORDER BY fm2.user_id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []int64
	for rows.Next() {
		var friendUserID int64
		if err := rows.Scan(&friendUserID); err != nil {
			return nil, err
		}
		result = append(result, friendUserID)
	}

	return result, rows.Err()
}

func (r *FriendshipRepository) FindByUserPair(ctx context.Context, userAID, userBID int64) (*models.Friendship, error) {
	query := `
		SELECT f.id, f.status, f.created_at, f.updated_at, f.ended_at
		FROM friendships f
		JOIN friendship_members fm1 ON fm1.friendship_id = f.id
		JOIN friendship_members fm2 ON fm2.friendship_id = f.id
		WHERE fm1.user_id = $1
		  AND fm2.user_id = $2
		  AND fm1.user_id <> fm2.user_id
		LIMIT 1
	`

	var fs models.Friendship
	err := r.db.QueryRowContext(ctx, query, userAID, userBID).Scan(
		&fs.ID,
		&fs.Status,
		&fs.CreatedAt,
		&fs.UpdatedAt,
		&fs.EndedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &fs, nil
}

func (r *FriendshipRepository) CountByUserID(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COUNT(DISTINCT f.id)
		FROM friendships f
		JOIN friendship_members fm ON fm.friendship_id = f.id
		WHERE fm.user_id = $1
	`

	var total int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *FriendshipRepository) UpdateStatusTx(
	ctx context.Context,
	tx *sql.Tx,
	friendshipID int64,
	status string,
	endedAt sql.NullTime,
) error {
	query := `
		UPDATE friendships
		SET status = $1,
			ended_at = $2,
			updated_at = $3
		WHERE id = $4
	`

	_, err := tx.ExecContext(ctx, query, status, endedAt, time.Now(), friendshipID)
	return err
}
func (r *FriendshipRepository) GetOtherUserID(
	ctx context.Context,
	friendshipID int64,
	currentUserID int64,
) (*int64, error) {
	query := `
		SELECT user_id
		FROM friendship_members
		WHERE friendship_id = $1
		  AND user_id <> $2
		LIMIT 1
	`

	var otherUserID int64
	err := r.db.QueryRowContext(ctx, query, friendshipID, currentUserID).Scan(&otherUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &otherUserID, nil
}

func (r *FriendshipRepository) ExistsActivePair(
	ctx context.Context,
	userAID, userBID int64,
) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM friendships f
			JOIN friendship_members fm1 ON fm1.friendship_id = f.id
			JOIN friendship_members fm2 ON fm2.friendship_id = f.id
			WHERE f.status = 'active'
			  AND fm1.user_id = $1
			  AND fm2.user_id = $2
			  AND fm1.user_id <> fm2.user_id
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userAID, userBID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
func (r *FriendshipRepository) ListMutualFriendUserIDs(
	ctx context.Context,
	userAID, userBID int64,
) ([]int64, error) {
	query := `
		SELECT DISTINCT a.friend_user_id
		FROM (
			SELECT fm2.user_id AS friend_user_id
			FROM friendships f
			JOIN friendship_members fm1 ON fm1.friendship_id = f.id
			JOIN friendship_members fm2 ON fm2.friendship_id = f.id
			WHERE f.status = 'active'
			  AND fm1.user_id = $1
			  AND fm2.user_id <> $1
		) a
		JOIN (
			SELECT fm2.user_id AS friend_user_id
			FROM friendships f
			JOIN friendship_members fm1 ON fm1.friendship_id = f.id
			JOIN friendship_members fm2 ON fm2.friendship_id = f.id
			WHERE f.status = 'active'
			  AND fm1.user_id = $2
			  AND fm2.user_id <> $2
		) b ON a.friend_user_id = b.friend_user_id
		ORDER BY a.friend_user_id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userAID, userBID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []int64
	for rows.Next() {
		var friendUserID int64
		if err := rows.Scan(&friendUserID); err != nil {
			return nil, err
		}
		result = append(result, friendUserID)
	}

	return result, rows.Err()
}