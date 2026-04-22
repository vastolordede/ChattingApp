package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type FriendRequestRepository struct {
	db *sql.DB
}

func NewFriendRequestRepository(db *sql.DB) *FriendRequestRepository {
	return &FriendRequestRepository{db: db}
}

func (r *FriendRequestRepository) CreateTx(ctx context.Context, tx *sql.Tx, fr *models.FriendRequest) (int64, error) {
	query := `
		INSERT INTO friend_requests (
			status, message, created_at, updated_at, responded_at, expired_at
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int64
	err := tx.QueryRowContext(
		ctx,
		query,
		fr.Status,
		fr.Message,
		fr.CreatedAt,
		fr.UpdatedAt,
		fr.RespondedAt,
		fr.ExpiredAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *FriendRequestRepository) AddMemberTx(ctx context.Context, tx *sql.Tx, member *models.FriendRequestMember) error {
	query := `
		INSERT INTO friend_request_members (
			friend_request_id, user_id, role, created_at
		) VALUES ($1, $2, $3, $4)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		member.FriendRequestID,
		member.UserID,
		member.Role,
		member.CreatedAt,
	)

	return err
}

func (r *FriendRequestRepository) FindPendingBetweenUsers(ctx context.Context, userAID, userBID int64) (*models.FriendRequest, error) {
	query := `
		SELECT fr.id, fr.status, fr.message, fr.created_at, fr.updated_at, fr.responded_at, fr.expired_at
		FROM friend_requests fr
		JOIN friend_request_members frm1 ON frm1.friend_request_id = fr.id
		JOIN friend_request_members frm2 ON frm2.friend_request_id = fr.id
		WHERE fr.status = 'pending'
		  AND frm1.user_id = $1
		  AND frm2.user_id = $2
		  AND frm1.role = 'sender'
		  AND frm2.role = 'receiver'
		  AND frm1.user_id <> frm2.user_id
		LIMIT 1
	`

	var fr models.FriendRequest
	err := r.db.QueryRowContext(ctx, query, userAID, userBID).Scan(
		&fr.ID,
		&fr.Status,
		&fr.Message,
		&fr.CreatedAt,
		&fr.UpdatedAt,
		&fr.RespondedAt,
		&fr.ExpiredAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &fr, nil
}

func (r *FriendRequestRepository) GetByID(ctx context.Context, id int64) (*models.FriendRequest, error) {
	query := `
		SELECT id, status, message, created_at, updated_at, responded_at, expired_at
		FROM friend_requests
		WHERE id = $1
		LIMIT 1
	`

	var fr models.FriendRequest
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&fr.ID,
		&fr.Status,
		&fr.Message,
		&fr.CreatedAt,
		&fr.UpdatedAt,
		&fr.RespondedAt,
		&fr.ExpiredAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &fr, nil
}

func (r *FriendRequestRepository) GetMembersByRequestID(ctx context.Context, requestID int64) ([]models.FriendRequestMember, error) {
	query := `
		SELECT id, friend_request_id, user_id, role, created_at
		FROM friend_request_members
		WHERE friend_request_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.FriendRequestMember
	for rows.Next() {
		var item models.FriendRequestMember
		if err := rows.Scan(
			&item.ID,
			&item.FriendRequestID,
			&item.UserID,
			&item.Role,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func (r *FriendRequestRepository) UpdateStatusTx(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
	status string,
	respondedAt sql.NullTime,
) error {
	query := `
		UPDATE friend_requests
		SET status = $1,
			responded_at = $2,
			updated_at = $3
		WHERE id = $4
	`

	_, err := tx.ExecContext(ctx, query, status, respondedAt, time.Now(), id)
	return err
}

func (r *FriendRequestRepository) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]models.FriendRequest, error) {
	query := `
		SELECT DISTINCT fr.id, fr.status, fr.message, fr.created_at, fr.updated_at, fr.responded_at, fr.expired_at
		FROM friend_requests fr
		JOIN friend_request_members frm ON frm.friend_request_id = fr.id
		WHERE frm.user_id = $1
		ORDER BY fr.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.FriendRequest
	for rows.Next() {
		var item models.FriendRequest
		if err := rows.Scan(
			&item.ID,
			&item.Status,
			&item.Message,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.RespondedAt,
			&item.ExpiredAt,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func (r *FriendRequestRepository) CountByUserID(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COUNT(DISTINCT fr.id)
		FROM friend_requests fr
		JOIN friend_request_members frm ON frm.friend_request_id = fr.id
		WHERE frm.user_id = $1
	`

	var total int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
func (r *FriendRequestRepository) GetSenderAndReceiverByRequestID(
	ctx context.Context,
	requestID int64,
) (senderID int64, receiverID int64, err error) {
	query := `
		SELECT user_id, role
		FROM friend_request_members
		WHERE friend_request_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, requestID)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int64
		var role string
		if err := rows.Scan(&userID, &role); err != nil {
			return 0, 0, err
		}

		switch role {
		case "sender":
			senderID = userID
		case "receiver":
			receiverID = userID
		}
	}

	if err := rows.Err(); err != nil {
		return 0, 0, err
	}

	return senderID, receiverID, nil
}

func (r *FriendRequestRepository) ListIncomingPendingByUserID(
	ctx context.Context,
	userID int64,
	limit, offset int,
) ([]models.FriendRequest, error) {
	query := `
		SELECT fr.id, fr.status, fr.message, fr.created_at, fr.updated_at, fr.responded_at, fr.expired_at
		FROM friend_requests fr
		JOIN friend_request_members sender_m
			ON sender_m.friend_request_id = fr.id AND sender_m.role = 'sender'
		JOIN friend_request_members receiver_m
			ON receiver_m.friend_request_id = fr.id AND receiver_m.role = 'receiver'
		WHERE fr.status = 'pending'
		  AND receiver_m.user_id = $1
		ORDER BY fr.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.FriendRequest
	for rows.Next() {
		var item models.FriendRequest
		if err := rows.Scan(
			&item.ID,
			&item.Status,
			&item.Message,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.RespondedAt,
			&item.ExpiredAt,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}