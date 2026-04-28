package service

import (
	"chattingapp_be/internal/dto"
	"chattingapp_be/internal/models"
	"chattingapp_be/internal/repository"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type FriendService struct {
	db             *sql.DB
	friendReqRepo  *repository.FriendRequestRepository
	friendshipRepo *repository.FriendshipRepository
	userRepo       *repository.UserRepository
}

func NewFriendService(
	db *sql.DB,
	friendReqRepo *repository.FriendRequestRepository,
	friendshipRepo *repository.FriendshipRepository,
	userRepo *repository.UserRepository,
) *FriendService {
	return &FriendService{
		db:             db,
		friendReqRepo:  friendReqRepo,
		friendshipRepo: friendshipRepo,
		userRepo:       userRepo,
	}
}

func toFriendUserSummary(u *models.User) dto.FriendUserSummary {
	var avatarURL *string
	if u.AvatarURL.Valid {
		avatarURL = &u.AvatarURL.String
	}

	return dto.FriendUserSummary{
		ID:         u.ID,
		Username:   u.Username,
		FullName:   u.FullName,
		AvatarURL:  avatarURL,
		IsVerified: u.IsVerified,
	}
}

func toOptionalString(t sql.NullTime) *string {
	if !t.Valid {
		return nil
	}
	s := t.Time.Format(time.RFC3339)
	return &s
}

func (s *FriendService) buildFriendRequestResponse(
	ctx context.Context,
	fr models.FriendRequest,
) (*dto.FriendRequestResponse, error) {
	senderID, receiverID, err := s.friendReqRepo.GetSenderAndReceiverByRequestID(ctx, fr.ID)
	if err != nil {
		return nil, err
	}
	if senderID == 0 || receiverID == 0 {
		return nil, errors.New("friend request members không hợp lệ")
	}

	sender, err := s.userRepo.GetByID(ctx, senderID)
	if err != nil {
		return nil, err
	}
	if sender == nil {
		return nil, errors.New("sender không tồn tại")
	}

	receiver, err := s.userRepo.GetByID(ctx, receiverID)
	if err != nil {
		return nil, err
	}
	if receiver == nil {
		return nil, errors.New("receiver không tồn tại")
	}

	var message *string
	if fr.Message.Valid {
		message = &fr.Message.String
	}

	return &dto.FriendRequestResponse{
		ID:          fr.ID,
		Status:      fr.Status,
		Message:     message,
		CreatedAt:   fr.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   fr.UpdatedAt.Format(time.RFC3339),
		RespondedAt: toOptionalString(fr.RespondedAt),
		ExpiredAt:   toOptionalString(fr.ExpiredAt),
		Sender:      toFriendUserSummary(sender),
		Receiver:    toFriendUserSummary(receiver),
	}, nil
}
func (s *FriendService) SendFriendRequest(
	ctx context.Context,
	senderUserID int64,
	req dto.SendFriendRequestRequest,
) error {
	if senderUserID == req.ReceiverUserID {
		return errors.New("không thể gửi lời mời kết bạn cho chính mình")
	}

	receiver, err := s.userRepo.GetByID(ctx, req.ReceiverUserID)
	if err != nil {
		return err
	}
	if receiver == nil {
		return errors.New("người nhận không tồn tại")
	}

	existingFriendship, err := s.friendshipRepo.FindActiveByUserPair(ctx, senderUserID, req.ReceiverUserID)
	if err != nil {
		return err
	}
	if existingFriendship != nil {
		return errors.New("hai người đã là bạn bè")
	}

	pendingA, err := s.friendReqRepo.FindPendingBetweenUsers(ctx, senderUserID, req.ReceiverUserID)
	if err != nil {
		return err
	}
	if pendingA != nil {
		return errors.New("đã tồn tại lời mời kết bạn đang chờ")
	}

	pendingB, err := s.friendReqRepo.FindPendingBetweenUsers(ctx, req.ReceiverUserID, senderUserID)
	if err != nil {
		return err
	}
	if pendingB != nil {
		return errors.New("đối phương đã gửi lời mời cho bạn")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	fr := &models.FriendRequest{
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if req.Message != nil {
		fr.Message = sql.NullString{String: *req.Message, Valid: true}
	}

	friendRequestID, err := s.friendReqRepo.CreateTx(ctx, tx, fr)
	if err != nil {
		return err
	}

	if err := s.friendReqRepo.AddMemberTx(ctx, tx, &models.FriendRequestMember{
		FriendRequestID: friendRequestID,
		UserID:          senderUserID,
		Role:            "sender",
		CreatedAt:       now,
	}); err != nil {
		return err
	}

	if err := s.friendReqRepo.AddMemberTx(ctx, tx, &models.FriendRequestMember{
		FriendRequestID: friendRequestID,
		UserID:          req.ReceiverUserID,
		Role:            "receiver",
		CreatedAt:       now,
	}); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *FriendService) RespondFriendRequest(
	ctx context.Context,
	currentUserID int64,
	friendRequestID int64,
	req dto.RespondFriendRequestRequest,
) error {
	fr, err := s.friendReqRepo.GetByID(ctx, friendRequestID)
	if err != nil {
		return err
	}
	if fr == nil {
		return errors.New("không tìm thấy lời mời kết bạn")
	}
	if fr.Status != "pending" {
		return errors.New("lời mời này không còn ở trạng thái pending")
	}

	members, err := s.friendReqRepo.GetMembersByRequestID(ctx, friendRequestID)
	if err != nil {
		return err
	}

	var senderID int64
	var receiverID int64
	for _, m := range members {
		if m.Role == "sender" {
			senderID = m.UserID
		}
		if m.Role == "receiver" {
			receiverID = m.UserID
		}
	}

	if senderID == 0 || receiverID == 0 {
		return errors.New("dữ liệu member của friend request không hợp lệ")
	}

	switch req.Action {
	case "accepted":
		if currentUserID != receiverID {
			return errors.New("chỉ người nhận mới được chấp nhận lời mời")
		}

		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		now := time.Now()
		if err := s.friendReqRepo.UpdateStatusTx(
			ctx,
			tx,
			friendRequestID,
			"accepted",
			sql.NullTime{Time: now, Valid: true},
		); err != nil {
			return err
		}

		friendshipID, err := s.friendshipRepo.CreateTx(ctx, tx, &models.Friendship{
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			return err
		}

		if err := s.friendshipRepo.AddMemberTx(ctx, tx, &models.FriendshipMember{
			FriendshipID: friendshipID,
			UserID:       senderID,
			CreatedAt:    now,
		}); err != nil {
			return err
		}

		if err := s.friendshipRepo.AddMemberTx(ctx, tx, &models.FriendshipMember{
			FriendshipID: friendshipID,
			UserID:       receiverID,
			CreatedAt:    now,
		}); err != nil {
			return err
		}

		return tx.Commit()

	case "rejected":
		if currentUserID != receiverID {
			return errors.New("chỉ người nhận mới được từ chối lời mời")
		}

		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		if err := s.friendReqRepo.UpdateStatusTx(
			ctx,
			tx,
			friendRequestID,
			"rejected",
			sql.NullTime{Time: time.Now(), Valid: true},
		); err != nil {
			return err
		}

		return tx.Commit()

	case "cancelled":
		if currentUserID != senderID {
			return errors.New("chỉ người gửi mới được hủy lời mời")
		}

		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		if err := s.friendReqRepo.UpdateStatusTx(
			ctx,
			tx,
			friendRequestID,
			"cancelled",
			sql.NullTime{Time: time.Now(), Valid: true},
		); err != nil {
			return err
		}

		return tx.Commit()

	default:
		return errors.New("action không hợp lệ")
	}
}

func (s *FriendService) ListFriends(ctx context.Context, userID int64) ([]dto.FriendUserSummary, error) {
	friendIDs, err := s.friendshipRepo.ListFriendUserIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.FriendUserSummary, 0, len(friendIDs))
	for _, fid := range friendIDs {
		u, err := s.userRepo.GetByID(ctx, fid)
		if err != nil {
			return nil, err
		}
		if u == nil {
			continue
		}

		var avatarURL *string
		if u.AvatarURL.Valid {
			avatarURL = &u.AvatarURL.String
		}

		item := dto.FriendUserSummary{
			ID:         u.ID,
			Username:   u.Username,
			FullName:   u.FullName,
			AvatarURL:  avatarURL,
			IsVerified: u.IsVerified,
		}
		result = append(result, item)
	}

	return result, nil
}
func (s *FriendService) ListIncomingRequests(ctx context.Context, userID int64) ([]dto.FriendRequestResponse, error) {
	friendRequests, err := s.friendReqRepo.ListIncomingPendingByUserID(ctx, userID, 50, 0)
	if err != nil {
		return nil, err
	}

	result := make([]dto.FriendRequestResponse, 0, len(friendRequests))
	for _, fr := range friendRequests {
		item, err := s.buildFriendRequestResponse(ctx, fr)
		if err != nil {
			return nil, err
		}
		result = append(result, *item)
	}

	return result, nil
}

func (s *FriendService) ListOutgoingRequests(ctx context.Context, userID int64) ([]dto.FriendRequestResponse, error) {
	friendRequests, err := s.friendReqRepo.ListOutgoingPendingByUserID(ctx, userID, 50, 0)
	if err != nil {
		return nil, err
	}

	result := make([]dto.FriendRequestResponse, 0, len(friendRequests))
	for _, fr := range friendRequests {
		item, err := s.buildFriendRequestResponse(ctx, fr)
		if err != nil {
			return nil, err
		}
		result = append(result, *item)
	}

	return result, nil
}

func (s *FriendService) CancelFriendRequest(ctx context.Context, currentUserID, friendRequestID int64) error {
	return s.RespondFriendRequest(ctx, currentUserID, friendRequestID, dto.RespondFriendRequestRequest{
		Action: "cancelled",
	})
}

func (s *FriendService) Unfriend(ctx context.Context, currentUserID, targetUserID int64) error {
	if currentUserID == targetUserID {
		return errors.New("không thể unfriend chính mình")
	}

	friendship, err := s.friendshipRepo.FindActiveByUserPair(ctx, currentUserID, targetUserID)
	if err != nil {
		return err
	}
	if friendship == nil {
		return errors.New("hai người chưa là bạn bè")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.friendshipRepo.UpdateStatusTx(
		ctx,
		tx,
		friendship.ID,
		"ended",
		sql.NullTime{Time: time.Now(), Valid: true},
	); err != nil {
		return err
	}

	return tx.Commit()
}
func (s *FriendService) SearchUsers(ctx context.Context, currentUserID int64, q string) ([]dto.FriendUserSummary, error) {
	q = strings.TrimSpace(q)
	if q == "" {
		return []dto.FriendUserSummary{}, nil
	}

	query := `
		SELECT id, username, full_name, avatar_url, is_verified
		FROM users
		WHERE id <> $1
		  AND (
			username ILIKE '%' || $2 || '%'
			OR full_name ILIKE '%' || $2 || '%'
		  )
		ORDER BY is_verified DESC, username ASC
		LIMIT 20
	`

	rows, err := s.db.QueryContext(ctx, query, currentUserID, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.FriendUserSummary
	for rows.Next() {
		var item dto.FriendUserSummary
		var avatarURL sql.NullString

		if err := rows.Scan(
			&item.ID,
			&item.Username,
			&item.FullName,
			&avatarURL,
			&item.IsVerified,
		); err != nil {
			return nil, err
		}

		if avatarURL.Valid {
			item.AvatarURL = &avatarURL.String
		}

		result = append(result, item)
	}

	return result, rows.Err()
}

func (s *FriendService) ListMutualFriends(ctx context.Context, currentUserID, targetUserID int64) ([]dto.FriendUserSummary, error) {
	if currentUserID == targetUserID {
		return nil, errors.New("không thể lấy mutual friends với chính mình")
	}

	targetUser, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}
	if targetUser == nil {
		return nil, errors.New("user không tồn tại")
	}

	friendIDs, err := s.friendshipRepo.ListMutualFriendUserIDs(ctx, currentUserID, targetUserID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.FriendUserSummary, 0, len(friendIDs))
	for _, fid := range friendIDs {
		u, err := s.userRepo.GetByID(ctx, fid)
		if err != nil {
			return nil, err
		}
		if u == nil {
			continue
		}

		item := toFriendUserSummary(u)
		result = append(result, item)
	}

	return result, nil
}

func (s *FriendService) BlockUser(ctx context.Context, currentUserID, targetUserID int64) error {
	if currentUserID == targetUserID {
		return errors.New("không thể block chính mình")
	}

	targetUser, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return err
	}
	if targetUser == nil {
		return errors.New("user không tồn tại")
	}

	var exists bool
	err = s.db.QueryRowContext(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM user_blocks
			WHERE blocker_user_id = $1 AND blocked_user_id = $2
		)`,
		currentUserID,
		targetUserID,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user đã bị block trước đó")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Nếu đang là bạn thì kết thúc friendship
	friendship, err := s.friendshipRepo.FindActiveByUserPair(ctx, currentUserID, targetUserID)
	if err != nil {
		return err
	}
	if friendship != nil {
		if err := s.friendshipRepo.UpdateStatusTx(
			ctx,
			tx,
			friendship.ID,
			"ended",
			sql.NullTime{Time: time.Now(), Valid: true},
		); err != nil {
			return err
		}
	}

	// Đóng mọi lời mời pending giữa 2 bên
	_, err = tx.ExecContext(ctx, `
		UPDATE friend_requests fr
		SET status = 'cancelled',
			responded_at = NOW(),
			updated_at = NOW()
		WHERE fr.status = 'pending'
		  AND fr.id IN (
			SELECT fr2.id
			FROM friend_requests fr2
			JOIN friend_request_members sender_m
				ON sender_m.friend_request_id = fr2.id AND sender_m.role = 'sender'
			JOIN friend_request_members receiver_m
				ON receiver_m.friend_request_id = fr2.id AND receiver_m.role = 'receiver'
			WHERE (sender_m.user_id = $1 AND receiver_m.user_id = $2)
			   OR (sender_m.user_id = $2 AND receiver_m.user_id = $1)
		  )
	`, currentUserID, targetUserID)
	if err != nil {
		return err
	}

	// Tạo block
	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_blocks (blocker_user_id, blocked_user_id, created_at)
		VALUES ($1, $2, NOW())
	`, currentUserID, targetUserID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
