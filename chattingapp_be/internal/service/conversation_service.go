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
	"chattingapp_be/internal/realtime"
)

type ConversationService struct {
	db                     *sql.DB
	conversationRepo       *repository.ConversationRepository
	conversationMemberRepo *repository.ConversationMemberRepository
	userRepo               *repository.UserRepository
	messageRepo            *repository.MessageRepository
	realtimeHub *realtime.Hub
}

func NewConversationService(
	db *sql.DB,
	conversationRepo *repository.ConversationRepository,
	conversationMemberRepo *repository.ConversationMemberRepository,
	userRepo *repository.UserRepository,
	messageRepo *repository.MessageRepository,
	realtimeHub *realtime.Hub,
) *ConversationService {
	return &ConversationService{
		db:                     db,
		conversationRepo:       conversationRepo,
		conversationMemberRepo: conversationMemberRepo,
		userRepo:               userRepo,
		messageRepo:            messageRepo,
		realtimeHub: realtimeHub,
	}
}

func (s *ConversationService) CreateDirectConversation(
	ctx context.Context,
	userID int64,
	req dto.CreateDirectConversationRequest,
) (*dto.ConversationDetailResponse, error) {
	if userID == req.TargetUserID {
		return nil, errors.New("không thể tạo cuộc trò chuyện với chính mình")
	}

	targetUser, err := s.userRepo.GetByID(ctx, req.TargetUserID)
	if err != nil {
		return nil, err
	}
	if targetUser == nil {
		return nil, errors.New("user đích không tồn tại")
	}

	existing, err := s.conversationRepo.FindDirectConversationBetweenUsers(ctx, userID, req.TargetUserID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return s.GetConversationDetail(ctx, userID, existing.ID)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()

	conversation := &models.Conversation{
		ConversationType: "direct",
		Status:           "active",
		CreatedByUserID:  sql.NullInt64{Int64: userID, Valid: true},
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	conversationID, err := s.conversationRepo.CreateTx(ctx, tx, conversation)
	if err != nil {
		return nil, err
	}

	memberA := &models.ConversationMember{
		ConversationID: conversationID,
		UserID:         userID,
		Role:           "member",
		JoinedAt:       now,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	memberB := &models.ConversationMember{
		ConversationID: conversationID,
		UserID:         req.TargetUserID,
		Role:           "member",
		JoinedAt:       now,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if _, err := s.conversationMemberRepo.CreateTx(ctx, tx, memberA); err != nil {
		return nil, err
	}
	if _, err := s.conversationMemberRepo.CreateTx(ctx, tx, memberB); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetConversationDetail(ctx, userID, conversationID)
}

func (s *ConversationService) ListMyConversations(
	ctx context.Context,
	userID int64,
	page, limit int,
) ([]dto.ConversationListItemResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	offset := (page - 1) * limit

	conversations, err := s.conversationRepo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ConversationListItemResponse, 0, len(conversations))
	for _, c := range conversations {
		item := dto.ConversationListItemResponse{
			ID:               c.ID,
			ConversationType: c.ConversationType,
			Title:            nullStringToPtr(c.Title),
			AvatarURL:        nullStringToPtr(c.AvatarURL),
			Status:           c.Status,
			LastMessageID:    nullInt64ToPtr(c.LastMessageID),
			LastMessageAt:    nullTimeToPtrString(c.LastMessageAt),
			CreatedAt:        c.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        c.UpdatedAt.Format(time.RFC3339),
		}

		if c.LastMessageID.Valid {
			lastMsg, err := s.messageRepo.GetByID(ctx, c.LastMessageID.Int64)
			if err != nil {
				return nil, err
			}
			if lastMsg != nil {
				item.LastMessage = toMessagePreviewResponse(lastMsg)
			}
		}

		result = append(result, item)
	}

	return result, nil
}

func (s *ConversationService) GetConversationDetail(
	ctx context.Context,
	userID, conversationID int64,
) (*dto.ConversationDetailResponse, error) {
	member, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil || !member.IsActive {
		return nil, errors.New("bạn không thuộc cuộc trò chuyện này")
	}

	conversation, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	if conversation == nil {
		return nil, errors.New("không tìm thấy cuộc trò chuyện")
	}

	members, err := s.conversationMemberRepo.ListByConversationID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	memberResponses := make([]dto.ConversationMemberSummary, 0, len(members))
	for _, m := range members {
		u, err := s.userRepo.GetByID(ctx, m.UserID)
		if err != nil {
			return nil, err
		}
		if u == nil {
			continue
		}

		memberResponses = append(memberResponses, dto.ConversationMemberSummary{
			UserID:      u.ID,
			Username:    u.Username,
			FullName:    u.FullName,
			AvatarURL:   nullStringToPtr(u.AvatarURL),
			Role:        m.Role,
			Nickname:    nullStringToPtr(m.Nickname),
			IsActive:    m.IsActive,
			JoinedAt:    m.JoinedAt.Format(time.RFC3339),
			LastReadAt:  nullTimeToPtrString(m.LastReadAt),
			LastReadMsg: nullInt64ToPtr(m.LastReadMessageID),
		})
	}

	resp := &dto.ConversationDetailResponse{
		ID:               conversation.ID,
		ConversationType: conversation.ConversationType,
		Title:            nullStringToPtr(conversation.Title),
		AvatarURL:        nullStringToPtr(conversation.AvatarURL),
		Status:           conversation.Status,
		CreatedByUserID:  nullInt64ToPtr(conversation.CreatedByUserID),
		LastMessageID:    nullInt64ToPtr(conversation.LastMessageID),
		LastMessageAt:    nullTimeToPtrString(conversation.LastMessageAt),
		CreatedAt:        conversation.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        conversation.UpdatedAt.Format(time.RFC3339),
		Members:          memberResponses,
	}

	return resp, nil
}

func (s *ConversationService) MarkConversationRead(
	ctx context.Context,
	userID, conversationID int64,
	req dto.MarkConversationReadRequest,
) error {
	_, err := s.requireMember(ctx, userID, conversationID)
	if err != nil {
		return err
	}

	msg, err := s.messageRepo.GetByID(ctx, req.LastReadMessageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("message không tồn tại")
	}
	if msg.ConversationID != conversationID {
		return errors.New("message không thuộc cuộc trò chuyện này")
	}

	if err := s.conversationMemberRepo.UpdateLastRead(ctx, conversationID, userID, req.LastReadMessageID); err != nil {
	return err
}

s.broadcastConversationEvent(ctx, conversationID, realtime.Event{
	Type:           "message_read",
	ConversationID: conversationID,
	UserID:         userID,
	MessageID:      req.LastReadMessageID,
})

return nil
}

func (s *ConversationService) UpdateMyNickname(
	ctx context.Context,
	userID, conversationID int64,
	req dto.UpdateConversationNicknameRequest,
) error {
	member, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, conversationID, userID)
	if err != nil {
		return err
	}
	if member == nil || !member.IsActive {
		return errors.New("bạn không thuộc cuộc trò chuyện này")
	}

	var nickname sql.NullString
	if req.Nickname != nil {
		v := strings.TrimSpace(*req.Nickname)
		if v != "" {
			nickname = sql.NullString{String: v, Valid: true}
		}
	}

	return s.conversationMemberRepo.UpdateNickname(ctx, conversationID, userID, nickname)
}

func (s *ConversationService) MuteConversation(
	ctx context.Context,
	userID, conversationID int64,
	req dto.MuteConversationRequest,
) error {
	_, err := s.requireMember(ctx, userID, conversationID)
	if err != nil {
		return err
	}

	var muteUntil sql.NullTime
	if req.MuteUntil != nil {
		t, err := time.Parse(time.RFC3339, *req.MuteUntil)
		if err != nil {
			return errors.New("mute_until phải theo định dạng RFC3339")
		}
		if t.Before(time.Now()) {
			return errors.New("mute_until phải lớn hơn thời điểm hiện tại")
		}
		muteUntil = sql.NullTime{Time: t, Valid: true}
	}

	return s.conversationMemberRepo.UpdateMuteUntil(ctx, conversationID, userID, muteUntil)
}

func toMessagePreviewResponse(m *models.Message) *dto.MessagePreviewResponse {
	if m == nil {
		return nil
	}

	return &dto.MessagePreviewResponse{
		ID:           m.ID,
		SenderUserID: m.SenderUserID,
		MessageType:  m.MessageType,
		Content:      nullStringToPtr(m.Content),
		Status:       m.Status,
		IsEdited:     m.IsEdited,
		IsDeleted:    m.IsDeleted,
		SentAt:       m.SentAt.Format(time.RFC3339),
	}
}

func nullStringToPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

func nullTimeToPtrString(v sql.NullTime) *string {
	if !v.Valid {
		return nil
	}
	s := v.Time.Format(time.RFC3339)
	return &s
}

func nullInt64ToPtr(v sql.NullInt64) *int64 {
	if !v.Valid {
		return nil
	}
	x := v.Int64
	return &x
}
func (s *ConversationService) requireMember(
	ctx context.Context,
	userID, conversationID int64,
) (*models.ConversationMember, error) {

	member, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil || !member.IsActive {
		return nil, errors.New("bạn không thuộc cuộc trò chuyện này")
	}

	return member, nil
}
func (s *ConversationService) PinConversation(
	ctx context.Context,
	conversationID, userID int64,
	isPinned bool,
) error {
	return s.conversationMemberRepo.UpdatePinned(ctx, conversationID, userID, isPinned)
}

func (s *ConversationService) ArchiveConversation(
	ctx context.Context,
	conversationID, userID int64,
	isArchived bool,
) error {
	return s.conversationMemberRepo.UpdateArchived(ctx, conversationID, userID, isArchived)
}

func (s *ConversationService) GetUnreadCount(
	ctx context.Context,
	userID int64,
) (int64, error) {
	return s.conversationMemberRepo.CountUnread(ctx, userID)
}
func (s *ConversationService) broadcastConversationEvent(
	ctx context.Context,
	conversationID int64,
	event realtime.Event,
) {
	if s.realtimeHub == nil {
		return
	}

	members, err := s.conversationMemberRepo.ListByConversationID(ctx, conversationID)
	if err != nil {
		return
	}

	userIDs := make([]int64, 0, len(members))
	for _, m := range members {
		if m.IsActive {
			userIDs = append(userIDs, m.UserID)
		}
	}

	s.realtimeHub.BroadcastToUsers(userIDs, event)
}
func (s *ConversationService) SendTypingEvent(
	ctx context.Context,
	userID, conversationID int64,
	req dto.TypingRequest,
) error {
	_, err := s.requireMember(ctx, userID, conversationID)
	if err != nil {
		return err
	}

	s.broadcastConversationEvent(ctx, conversationID, realtime.Event{
		Type:           "typing",
		ConversationID: conversationID,
		UserID:         userID,
		IsTyping:       &req.IsTyping,
	})

	return nil
}
