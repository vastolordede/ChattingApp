package service

import (
	"chattingapp_be/internal/dto"
	"chattingapp_be/internal/models"
	"chattingapp_be/internal/realtime"
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
	messageCiphertextRepo  *repository.MessageCiphertextRepository
	userDeviceRepo         *repository.UserDeviceRepository
	realtimeHub            *realtime.Hub
}

func NewMessageService(
	db *sql.DB,
	messageRepo *repository.MessageRepository,
	messageAttachmentRepo *repository.MessageAttachmentRepository,
	conversationRepo *repository.ConversationRepository,
	conversationMemberRepo *repository.ConversationMemberRepository,
	userRepo *repository.UserRepository,
	messageCiphertextRepo *repository.MessageCiphertextRepository,
	userDeviceRepo *repository.UserDeviceRepository,
	realtimeHub *realtime.Hub,
) *MessageService {
	return &MessageService{
		db:                     db,
		messageRepo:            messageRepo,
		messageAttachmentRepo:  messageAttachmentRepo,
		conversationRepo:       conversationRepo,
		conversationMemberRepo: conversationMemberRepo,
		userRepo:               userRepo,
		messageCiphertextRepo:  messageCiphertextRepo,
		userDeviceRepo:         userDeviceRepo,
		realtimeHub:            realtimeHub,
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
	for _, attachmentReq := range req.Attachments {
		attachmentType := strings.TrimSpace(attachmentReq.AttachmentType)
		fileName := strings.TrimSpace(attachmentReq.FileName)
		mimeType := strings.TrimSpace(attachmentReq.MimeType)
		fileURL := strings.TrimSpace(attachmentReq.FileURL)

		if attachmentType == "" {
			return nil, errors.New("attachment_type không được để trống")
		}
		if fileName == "" {
			return nil, errors.New("file_name không được để trống")
		}
		if mimeType == "" {
			return nil, errors.New("mime_type không được để trống")
		}
		if attachmentReq.FileSize <= 0 {
			return nil, errors.New("file_size không hợp lệ")
		}
		if fileURL == "" {
			return nil, errors.New("file_url không được để trống")
		}

		attachment := &models.MessageAttachment{
			MessageID:      messageID,
			AttachmentType: attachmentType,
			FileName:       fileName,
			MimeType:       mimeType,
			FileSize:       attachmentReq.FileSize,
			FileURL:        fileURL,
			CreatedAt:      now,
		}

		if attachmentReq.ThumbnailURL != nil && strings.TrimSpace(*attachmentReq.ThumbnailURL) != "" {
			attachment.ThumbnailURL = sql.NullString{
				String: strings.TrimSpace(*attachmentReq.ThumbnailURL),
				Valid:  true,
			}
		}
		if attachmentReq.Width != nil {
			attachment.Width = sql.NullInt64{
				Int64: int64(*attachmentReq.Width),
				Valid: true,
			}
		}
		if attachmentReq.Height != nil {
			attachment.Height = sql.NullInt64{
				Int64: int64(*attachmentReq.Height),
				Valid: true,
			}
		}
		if attachmentReq.DurationSeconds != nil {
			attachment.DurationSeconds = sql.NullInt64{
				Int64: int64(*attachmentReq.DurationSeconds),
				Valid: true,
			}
		}
		if attachmentReq.Checksum != nil && strings.TrimSpace(*attachmentReq.Checksum) != "" {
			attachment.Checksum = sql.NullString{
				String: strings.TrimSpace(*attachmentReq.Checksum),
				Valid:  true,
			}
		}
		if attachmentReq.EncryptionKeyHint != nil && strings.TrimSpace(*attachmentReq.EncryptionKeyHint) != "" {
			attachment.EncryptionKeyHint = sql.NullString{
				String: strings.TrimSpace(*attachmentReq.EncryptionKeyHint),
				Valid:  true,
			}
		}

		if _, err := s.messageAttachmentRepo.CreateTx(ctx, tx, attachment); err != nil {
			return nil, err
		}
	}
	if err := s.conversationRepo.UpdateLastMessageTx(ctx, tx, req.ConversationID, messageID, now); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	resp, err := s.buildMessageResponse(ctx, messageID)
	if err != nil {
		return nil, err
	}

	s.broadcastMessageEvent(ctx, req.ConversationID, realtime.Event{
		Type:           "message_created",
		ConversationID: req.ConversationID,
		UserID:         userID,
		MessageID:      messageID,
		Payload:        resp,
	})

	return resp, nil
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
func (s *MessageService) RecallMessage(
	ctx context.Context,
	userID, messageID int64,
) (*dto.MessageResponse, error) {
	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, errors.New("không tìm thấy message")
	}
	if msg.SenderUserID != userID {
		return nil, errors.New("bạn không có quyền thu hồi tin nhắn này")
	}
	if msg.IsDeleted {
		return nil, errors.New("không thể thu hồi tin nhắn đã xóa")
	}
	if msg.IsRecalled {
		return nil, errors.New("tin nhắn đã được thu hồi")
	}

	member, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, msg.ConversationID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil || !member.IsActive {
		return nil, errors.New("bạn không thuộc cuộc trò chuyện này")
	}

	if err := s.messageRepo.Recall(ctx, messageID); err != nil {
		return nil, err
	}

	return s.buildMessageResponse(ctx, messageID)
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
		IsRecalled:             m.IsRecalled,
		RecalledAt:             nullTimeToPtrString(m.RecalledAt),
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
func (s *MessageService) SearchMessages(
	ctx context.Context,
	userID, conversationID int64,
	keyword string,
	page, limit int,
) ([]dto.MessageResponse, error) {
	member, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil || !member.IsActive {
		return nil, errors.New("bạn không thuộc cuộc trò chuyện này")
	}

	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, errors.New("từ khóa tìm kiếm không được để trống")
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	messages, err := s.messageRepo.SearchByConversationID(ctx, conversationID, keyword, limit, offset)
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

func (s *MessageService) ListMessagesBeforeID(
	ctx context.Context,
	userID, conversationID, beforeID int64,
	limit int,
) ([]dto.MessageResponse, error) {
	member, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil || !member.IsActive {
		return nil, errors.New("bạn không thuộc cuộc trò chuyện này")
	}

	if beforeID <= 0 {
		return nil, errors.New("before_id không hợp lệ")
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	messages, err := s.messageRepo.ListBeforeIDByConversationID(ctx, conversationID, beforeID, limit)
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

func (s *MessageService) ForwardMessage(
	ctx context.Context,
	userID, sourceMessageID int64,
	req dto.ForwardMessageRequest,
) (*dto.MessageResponse, error) {
	sourceMsg, err := s.messageRepo.GetByID(ctx, sourceMessageID)
	if err != nil {
		return nil, err
	}
	if sourceMsg == nil {
		return nil, errors.New("không tìm thấy message gốc")
	}
	if sourceMsg.IsDeleted {
		return nil, errors.New("không thể forward tin nhắn đã xóa")
	}

	sourceMember, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, sourceMsg.ConversationID, userID)
	if err != nil {
		return nil, err
	}
	if sourceMember == nil || !sourceMember.IsActive {
		return nil, errors.New("bạn không có quyền xem message gốc")
	}

	targetMember, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, req.TargetConversationID, userID)
	if err != nil {
		return nil, err
	}
	if targetMember == nil || !targetMember.IsActive {
		return nil, errors.New("bạn không thuộc conversation đích")
	}

	targetConversation, err := s.conversationRepo.GetByID(ctx, req.TargetConversationID)
	if err != nil {
		return nil, err
	}
	if targetConversation == nil {
		return nil, errors.New("không tìm thấy conversation đích")
	}
	if targetConversation.Status != "active" {
		return nil, errors.New("conversation đích không hoạt động")
	}

	now := time.Now()

	content := sourceMsg.Content
	if req.Content != nil && strings.TrimSpace(*req.Content) != "" {
		content = sql.NullString{
			String: strings.TrimSpace(*req.Content),
			Valid:  true,
		}
	}

	msg := &models.Message{
		ConversationID:         req.TargetConversationID,
		SenderUserID:           userID,
		MessageType:            sourceMsg.MessageType,
		Content:                content,
		ForwardedFromMessageID: sql.NullInt64{Int64: sourceMsg.ID, Valid: true},
		Status:                 "sent",
		IsEdited:               false,
		IsDeleted:              false,
		SentAt:                 now,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	if req.ClientMessageID != nil && strings.TrimSpace(*req.ClientMessageID) != "" {
		msg.ClientMessageID = sql.NullString{
			String: strings.TrimSpace(*req.ClientMessageID),
			Valid:  true,
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	newMessageID, err := s.messageRepo.CreateTx(ctx, tx, msg)
	if err != nil {
		return nil, err
	}

	if err := s.conversationRepo.UpdateLastMessageTx(ctx, tx, req.TargetConversationID, newMessageID, now); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.buildMessageResponse(ctx, newMessageID)
}

func (s *MessageService) ReactMessage(
	ctx context.Context,
	userID, messageID int64,
	req dto.ReactMessageRequest,
) (*dto.MessageReactionResponse, error) {
	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, errors.New("không tìm thấy message")
	}
	if msg.IsDeleted {
		return nil, errors.New("không thể reaction tin nhắn đã xóa")
	}

	member, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, msg.ConversationID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil || !member.IsActive {
		return nil, errors.New("bạn không thuộc cuộc trò chuyện này")
	}

	reactionType := strings.TrimSpace(req.ReactionType)
	if reactionType == "" {
		return nil, errors.New("reaction_type không được để trống")
	}

	reaction, err := s.messageRepo.UpsertReaction(ctx, messageID, userID, reactionType)
	if err != nil {
		return nil, err
	}

	return &dto.MessageReactionResponse{
		ID:           reaction.ID,
		MessageID:    reaction.MessageID,
		UserID:       reaction.UserID,
		ReactionType: reaction.ReactionType,
		CreatedAt:    reaction.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    reaction.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *MessageService) DeleteReaction(
	ctx context.Context,
	userID, messageID int64,
) error {
	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("không tìm thấy message")
	}

	member, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, msg.ConversationID, userID)
	if err != nil {
		return err
	}
	if member == nil || !member.IsActive {
		return errors.New("bạn không thuộc cuộc trò chuyện này")
	}

	return s.messageRepo.DeleteReaction(ctx, messageID, userID)
}

func (s *MessageService) ListReactions(
	ctx context.Context,
	userID, messageID int64,
) ([]dto.MessageReactionResponse, error) {
	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, errors.New("không tìm thấy message")
	}

	member, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, msg.ConversationID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil || !member.IsActive {
		return nil, errors.New("bạn không thuộc cuộc trò chuyện này")
	}

	reactions, err := s.messageRepo.ListReactionsByMessageID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.MessageReactionResponse, 0, len(reactions))
	for _, reaction := range reactions {
		result = append(result, dto.MessageReactionResponse{
			ID:           reaction.ID,
			MessageID:    reaction.MessageID,
			UserID:       reaction.UserID,
			ReactionType: reaction.ReactionType,
			CreatedAt:    reaction.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    reaction.UpdatedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}
func (s *MessageService) broadcastMessageEvent(
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
func toEncryptedCiphertextResponse(c models.MessageCiphertext) dto.EncryptedCiphertextResponse {
	var senderDeviceID *int64
	if c.SenderDeviceID.Valid {
		senderDeviceID = &c.SenderDeviceID.Int64
	}

	var encryptionHeader *string
	if c.EncryptionHeader.Valid {
		encryptionHeader = &c.EncryptionHeader.String
	}

	var nonce *string
	if c.Nonce.Valid {
		nonce = &c.Nonce.String
	}

	var deliveredAt *string
	if c.DeliveredAt.Valid {
		s := c.DeliveredAt.Time.Format(time.RFC3339)
		deliveredAt = &s
	}

	return dto.EncryptedCiphertextResponse{
		ID:               c.ID,
		MessageID:        c.MessageID,
		TargetDeviceID:   c.TargetDeviceID,
		SenderDeviceID:   senderDeviceID,
		Ciphertext:       c.Ciphertext,
		EncryptionHeader: encryptionHeader,
		Nonce:            nonce,
		Algorithm:        c.Algorithm,
		MessageVersion:   c.MessageVersion,
		IsDelivered:      c.IsDelivered,
		DeliveredAt:      deliveredAt,
		CreatedAt:        c.CreatedAt.Format(time.RFC3339),
	}
}
func (s *MessageService) SendEncryptedMessage(
	ctx context.Context,
	senderUserID int64,
	req dto.SendEncryptedMessageRequest,
) (*dto.SendEncryptedMessageResponse, error) {
	if req.ConversationID <= 0 {
		return nil, errors.New("conversation_id is required")
	}
	if strings.TrimSpace(req.SenderDeviceUUID) == "" {
		return nil, errors.New("sender_device_uuid is required")
	}
	if len(req.Ciphertexts) == 0 {
		return nil, errors.New("ciphertexts is required")
	}

	senderMember, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, req.ConversationID, senderUserID)
	if err != nil {
		return nil, err
	}
	if senderMember == nil || !senderMember.IsActive {
		return nil, errors.New("sender is not active member of conversation")
	}

	senderDevice, err := s.userDeviceRepo.GetByUUIDAndUserID(ctx, strings.TrimSpace(req.SenderDeviceUUID), senderUserID)
	if err != nil {
		return nil, err
	}
	if senderDevice == nil || !senderDevice.IsActive {
		return nil, errors.New("sender device not found or inactive")
	}

	seenTargetDevices := make(map[int64]bool)
	for _, item := range req.Ciphertexts {
		if item.TargetDeviceID <= 0 {
			return nil, errors.New("target_device_id must be greater than 0")
		}
		if seenTargetDevices[item.TargetDeviceID] {
			return nil, errors.New("duplicated target_device_id in ciphertexts")
		}
		seenTargetDevices[item.TargetDeviceID] = true

		if strings.TrimSpace(item.Ciphertext) == "" {
			return nil, errors.New("ciphertext is required")
		}

		targetDevice, err := s.userDeviceRepo.GetByID(ctx, item.TargetDeviceID)
		if err != nil {
			return nil, err
		}
		if targetDevice == nil || !targetDevice.IsActive {
			return nil, errors.New("target device not found or inactive")
		}

		targetMember, err := s.conversationMemberRepo.GetByConversationAndUser(ctx, req.ConversationID, targetDevice.UserID)
		if err != nil {
			return nil, err
		}
		if targetMember == nil || !targetMember.IsActive {
			return nil, errors.New("target device owner is not active member of conversation")
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()

	var clientMessageID sql.NullString
	if req.ClientMessageID != nil && strings.TrimSpace(*req.ClientMessageID) != "" {
		clientMessageID = sql.NullString{
			String: strings.TrimSpace(*req.ClientMessageID),
			Valid:  true,
		}
	}

	var replyToMessageID sql.NullInt64
	if req.ReplyToMessageID != nil && *req.ReplyToMessageID > 0 {
		replyToMessageID = sql.NullInt64{
			Int64: *req.ReplyToMessageID,
			Valid: true,
		}
	}

	message := &models.Message{
		ConversationID:         req.ConversationID,
		SenderUserID:           senderUserID,
		MessageType:            "encrypted",
		Content:                sql.NullString{},
		ReplyToMessageID:       replyToMessageID,
		ForwardedFromMessageID: sql.NullInt64{},
		Status:                 "sent",
		IsEdited:               false,
		EditedAt:               sql.NullTime{},
		IsDeleted:              false,
		DeletedAt:              sql.NullTime{},
		ClientMessageID:        clientMessageID,
		SentAt:                 now,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	messageID, err := s.messageRepo.CreateTx(ctx, tx, message)
	if err != nil {
		return nil, err
	}
	message.ID = messageID

	cipherResponses := make([]dto.EncryptedCiphertextResponse, 0, len(req.Ciphertexts))

	for _, item := range req.Ciphertexts {
		algorithm := strings.TrimSpace(item.Algorithm)
		if algorithm == "" {
			algorithm = "XCHACHA20_POLY1305"
		}

		version := item.MessageVersion
		if version <= 0 {
			version = 1
		}

		var encryptionHeader sql.NullString
		if item.EncryptionHeader != nil && strings.TrimSpace(*item.EncryptionHeader) != "" {
			encryptionHeader = sql.NullString{
				String: strings.TrimSpace(*item.EncryptionHeader),
				Valid:  true,
			}
		}

		var nonce sql.NullString
		if item.Nonce != nil && strings.TrimSpace(*item.Nonce) != "" {
			nonce = sql.NullString{
				String: strings.TrimSpace(*item.Nonce),
				Valid:  true,
			}
		}

		cipher := &models.MessageCiphertext{
			MessageID:        messageID,
			TargetDeviceID:   item.TargetDeviceID,
			SenderDeviceID:   sql.NullInt64{Int64: senderDevice.ID, Valid: true},
			Ciphertext:       strings.TrimSpace(item.Ciphertext),
			EncryptionHeader: encryptionHeader,
			Nonce:            nonce,
			Algorithm:        algorithm,
			MessageVersion:   version,
			IsDelivered:      false,
			DeliveredAt:      sql.NullTime{},
			CreatedAt:        now,
		}

		cipherID, err := s.messageCiphertextRepo.CreateTx(ctx, tx, cipher)
		if err != nil {
			return nil, err
		}
		cipher.ID = cipherID

		cipherResponses = append(cipherResponses, toEncryptedCiphertextResponse(*cipher))
	}

	if err := s.conversationRepo.UpdateLastMessageTx(ctx, tx, req.ConversationID, messageID, now); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	messageResp, err := s.buildMessageResponse(ctx, messageID)
	if err != nil {
		return nil, err
	}
	messageResp.Content = nil

	s.broadcastMessageEvent(ctx, req.ConversationID, realtime.Event{
		Type:           "encrypted_message_created",
		ConversationID: req.ConversationID,
		UserID:         senderUserID,
		MessageID:      messageID,
		Payload: map[string]any{
			"message_id":      messageID,
			"conversation_id": req.ConversationID,
			"sender_user_id":  senderUserID,
			"message_type":    "encrypted",
			"status":          "sent",
			"sent_at":         now.Format(time.RFC3339),
		},
	})

	return &dto.SendEncryptedMessageResponse{
		Message:     *messageResp,
		Ciphertexts: cipherResponses,
	}, nil
}
func (s *MessageService) ListUndeliveredCiphertextsForDevice(
	ctx context.Context,
	userID int64,
	deviceUUID string,
) ([]dto.EncryptedCiphertextResponse, error) {
	deviceUUID = strings.TrimSpace(deviceUUID)
	if deviceUUID == "" {
		return nil, errors.New("device_uuid is required")
	}

	device, err := s.userDeviceRepo.GetByUUIDAndUserID(ctx, deviceUUID, userID)
	if err != nil {
		return nil, err
	}
	if device == nil || !device.IsActive {
		return nil, errors.New("device not found or inactive")
	}

	items, err := s.messageCiphertextRepo.ListUndeliveredByTargetDeviceID(ctx, device.ID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.EncryptedCiphertextResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, toEncryptedCiphertextResponse(item))
	}

	return resp, nil
}
func (s *MessageService) MarkCiphertextDelivered(
	ctx context.Context,
	userID int64,
	ciphertextID int64,
) error {
	if ciphertextID <= 0 {
		return errors.New("ciphertext id is invalid")
	}

	cipher, err := s.messageCiphertextRepo.GetByID(ctx, ciphertextID)
	if err != nil {
		return err
	}
	if cipher == nil {
		return errors.New("ciphertext not found")
	}

	targetDevice, err := s.userDeviceRepo.GetByID(ctx, cipher.TargetDeviceID)
	if err != nil {
		return err
	}
	if targetDevice == nil {
		return errors.New("target device not found")
	}
	if targetDevice.UserID != userID {
		return errors.New("you are not allowed to mark this ciphertext delivered")
	}

	return s.messageCiphertextRepo.MarkDelivered(ctx, ciphertextID)
}
