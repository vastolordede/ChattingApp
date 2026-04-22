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

type MessageService struct {
	db                     *sql.DB
	messageRepo            *repository.MessageRepository
	messageAttachmentRepo  *repository.MessageAttachmentRepository
	conversationRepo       *repository.ConversationRepository
	conversationMemberRepo *repository.ConversationMemberRepository
	userRepo               *repository.UserRepository
}

func NewMessageService(
	db *sql.DB,
	messageRepo *repository.MessageRepository,
	messageAttachmentRepo *repository.MessageAttachmentRepository,
	conversationRepo *repository.ConversationRepository,
	conversationMemberRepo *repository.ConversationMemberRepository,
	userRepo *repository.UserRepository,
) *MessageService {
	return &MessageService{
		db:                     db,
		messageRepo:            messageRepo,
		messageAttachmentRepo:  messageAttachmentRepo,
		conversationRepo:       conversationRepo,
		conversationMemberRepo: conversationMemberRepo,
		userRepo:               userRepo,
	}
}

func (s *MessageService) SendMessage(
	ctx context.Context,
	userID int64,
	req dto.SendMessageRequest,
) (*dto.MessageResponse, error) {
	member, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, req.ConversationID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil || !member.IsActive {
		return nil, errors.New("bạn không thuộc cuộc trò chuyện này")
	}

	conversation, err := s.conversationRepo.GetByID(ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}
	if conversation == nil {
		return nil, errors.New("không tìm thấy cuộc trò chuyện")
	}
	if conversation.Status != "active" {
		return nil, errors.New("cuộc trò chuyện không hoạt động")
	}

	if req.MessageType == "" {
		return nil, errors.New("message_type không được để trống")
	}

	if req.MessageType == "text" {
		if req.Content == nil || strings.TrimSpace(*req.Content) == "" {
			return nil, errors.New("nội dung tin nhắn không được để trống")
		}
	}

	if req.ReplyToMessageID != nil {
		replyMsg, err := s.messageRepo.GetByID(ctx, *req.ReplyToMessageID)
		if err != nil {
			return nil, err
		}
		if replyMsg == nil {
			return nil, errors.New("reply_to_message không tồn tại")
		}
		if replyMsg.ConversationID != req.ConversationID {
			return nil, errors.New("reply_to_message không thuộc conversation này")
		}
	}

	now := time.Now()

	msg := &models.Message{
		ConversationID: req.ConversationID,
		SenderUserID:   userID,
		MessageType:    req.MessageType,
		Status:         "sent",
		IsEdited:       false,
		IsDeleted:      false,
		SentAt:         now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if req.Content != nil && strings.TrimSpace(*req.Content) != "" {
		msg.Content = sql.NullString{String: strings.TrimSpace(*req.Content), Valid: true}
	}
	if req.ReplyToMessageID != nil {
		msg.ReplyToMessageID = sql.NullInt64{Int64: *req.ReplyToMessageID, Valid: true}
	}
	if req.ClientMessageID != nil && strings.TrimSpace(*req.ClientMessageID) != "" {
		msg.ClientMessageID = sql.NullString{String: strings.TrimSpace(*req.ClientMessageID), Valid: true}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	messageID, err := s.messageRepo.CreateTx(ctx, tx, msg)
	if err != nil {
		return nil, err
	}

	if err := s.conversationRepo.UpdateLastMessageTx(ctx, tx, req.ConversationID, messageID, now); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.buildMessageResponse(ctx, messageID)
}

func (s *MessageService) EditMessage(
	ctx context.Context,
	userID, messageID int64,
	req dto.EditMessageRequest,
) (*dto.MessageResponse, error) {
	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, errors.New("không tìm thấy message")
	}
	if msg.SenderUserID != userID {
		return nil, errors.New("bạn không có quyền sửa tin nhắn này")
	}
	if msg.IsDeleted {
		return nil, errors.New("không thể sửa tin nhắn đã xóa")
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errors.New("nội dung không được để trống")
	}

	if err := s.messageRepo.UpdateContent(ctx, messageID, content); err != nil {
		return nil, err
	}

	return s.buildMessageResponse(ctx, messageID)
}

func (s *MessageService) DeleteMessage(
	ctx context.Context,
	userID, messageID int64,
	req dto.DeleteMessageRequest,
) error {
	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("không tìm thấy message")
	}
	if msg.SenderUserID != userID {
		return errors.New("bạn không có quyền xóa tin nhắn này")
	}

	if req.SoftDelete {
		return s.messageRepo.SoftDelete(ctx, messageID)
	}

	return s.messageRepo.HardDelete(ctx, messageID)
}

func (s *MessageService) ListMessages(
	ctx context.Context,
	userID, conversationID int64,
	page, limit int,
) ([]dto.MessageResponse, error) {
	member, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil || !member.IsActive {
		return nil, errors.New("bạn không thuộc cuộc trò chuyện này")
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	messages, err := s.messageRepo.ListByConversationID(ctx, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]dto.MessageResponse, 0, len(messages))
	for _, m := range messages {
		item, err := s.buildMessageResponseFromModel(ctx, &m)
		if err != nil {
			return nil, err
		}
		result = append(result, *item)
	}

	return result, nil
}

func (s *MessageService) buildMessageResponse(ctx context.Context, messageID int64) (*dto.MessageResponse, error) {
	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, errors.New("không tìm thấy message")
	}

	return s.buildMessageResponseFromModel(ctx, msg)
}

func (s *MessageService) buildMessageResponseFromModel(ctx context.Context, m *models.Message) (*dto.MessageResponse, error) {
	var sender *dto.MessageSenderInfo
	u, err := s.userRepo.GetByID(ctx, m.SenderUserID)
	if err != nil {
		return nil, err
	}
	if u != nil {
		sender = &dto.MessageSenderInfo{
			ID:        u.ID,
			Username:  u.Username,
			FullName:  u.FullName,
			AvatarURL: nullStringToPtr(u.AvatarURL),
		}
	}

	var replyPreview *dto.MessagePreviewResponse
	if m.ReplyToMessageID.Valid {
		replyMsg, err := s.messageRepo.GetByID(ctx, m.ReplyToMessageID.Int64)
		if err != nil {
			return nil, err
		}
		if replyMsg != nil {
			replyPreview = toMessagePreviewResponse(replyMsg)
		}
	}

	attachments, err := s.messageAttachmentRepo.ListByMessageID(ctx, m.ID)
	if err != nil {
		return nil, err
	}

	attachmentResponses := make([]dto.MessageAttachmentResponse, 0, len(attachments))
	for _, a := range attachments {
		attachmentResponses = append(attachmentResponses, dto.MessageAttachmentResponse{
			ID:                a.ID,
			AttachmentType:    a.AttachmentType,
			FileName:          a.FileName,
			MimeType:          a.MimeType,
			FileSize:          a.FileSize,
			FileURL:           a.FileURL,
			ThumbnailURL:      nullStringToPtr(a.ThumbnailURL),
			Width:             nullInt64ToIntPtr(a.Width),
			Height:            nullInt64ToIntPtr(a.Height),
			DurationSeconds:   nullInt64ToIntPtr(a.DurationSeconds),
			Checksum:          nullStringToPtr(a.Checksum),
			EncryptionKeyHint: nullStringToPtr(a.EncryptionKeyHint),
			CreatedAt:         a.CreatedAt.Format(time.RFC3339),
		})
	}

	resp := &dto.MessageResponse{
		ID:                     m.ID,
		ConversationID:         m.ConversationID,
		SenderUserID:           m.SenderUserID,
		MessageType:            m.MessageType,
		Content:                nullStringToPtr(m.Content),
		ReplyToMessageID:       nullInt64ToPtr(m.ReplyToMessageID),
		ForwardedFromMessageID: nullInt64ToPtr(m.ForwardedFromMessageID),
		Status:                 m.Status,
		IsEdited:               m.IsEdited,
		EditedAt:               nullTimeToPtrString(m.EditedAt),
		IsDeleted:              m.IsDeleted,
		DeletedAt:              nullTimeToPtrString(m.DeletedAt),
		ClientMessageID:        nullStringToPtr(m.ClientMessageID),
		SentAt:                 m.SentAt.Format(time.RFC3339),
		CreatedAt:              m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:              m.UpdatedAt.Format(time.RFC3339),
		Sender:                 sender,
		ReplyToMessagePreview:  replyPreview,
		Attachments:            attachmentResponses,
	}

	return resp, nil
}

func nullInt64ToIntPtr(v sql.NullInt64) *int {
	if !v.Valid {
		return nil
	}
	x := int(v.Int64)
	return &x
}