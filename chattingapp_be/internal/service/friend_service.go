package service

import (
	"chattingapp_be/internal/dto"
	"chattingapp_be/internal/models"
	"chattingapp_be/internal/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type FriendService struct {
	db                *sql.DB
	friendReqRepo     *repository.FriendRequestRepository
	friendshipRepo    *repository.FriendshipRepository
	userRepo          *repository.UserRepository
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