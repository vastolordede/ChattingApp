package handler

import (
	"chattingapp_be/internal/dto"
	"chattingapp_be/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type ConversationHandler struct {
	conversationService *service.ConversationService
}

func NewConversationHandler(conversationService *service.ConversationService) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
	}
}

func (h *ConversationHandler) CreateDirectConversation(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	var req dto.CreateDirectConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.conversationService.CreateDirectConversation(r.Context(), userID, req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "tạo direct conversation thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "tạo direct conversation thành công",
		Data:    resp,
	})
}

func (h *ConversationHandler) ListMyConversations(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	page := 1
	limit := 20

	if v := r.URL.Query().Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 {
			limit = l
		}
	}

	items, err := h.conversationService.ListMyConversations(r.Context(), userID, page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "lấy danh sách cuộc trò chuyện thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "lấy danh sách cuộc trò chuyện thành công",
		Data:    items,
	})
}

func (h *ConversationHandler) GetConversationDetail(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	conversationIDStr := r.PathValue("id")
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "conversation id không hợp lệ",
		})
		return
	}

	resp, err := h.conversationService.GetConversationDetail(r.Context(), userID, conversationID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "lấy chi tiết cuộc trò chuyện thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "lấy chi tiết cuộc trò chuyện thành công",
		Data:    resp,
	})
}

func (h *ConversationHandler) MarkConversationRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	conversationIDStr := r.PathValue("id")
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "conversation id không hợp lệ",
		})
		return
	}

	var req dto.MarkConversationReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	if err := h.conversationService.MarkConversationRead(r.Context(), userID, conversationID, req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "đánh dấu đã đọc thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "đánh dấu đã đọc thành công",
	})
}

func (h *ConversationHandler) UpdateMyNickname(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	conversationIDStr := r.PathValue("id")
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "conversation id không hợp lệ",
		})
		return
	}

	var req dto.UpdateConversationNicknameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	if err := h.conversationService.UpdateMyNickname(r.Context(), userID, conversationID, req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "cập nhật nickname thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "cập nhật nickname thành công",
	})
}

func (h *ConversationHandler) MuteConversation(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	conversationIDStr := r.PathValue("id")
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "conversation id không hợp lệ",
		})
		return
	}

	var req dto.MuteConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	if err := h.conversationService.MuteConversation(r.Context(), userID, conversationID, req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "mute conversation thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "mute conversation thành công",
	})
}
func (h *ConversationHandler) PinConversation(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	idStr := r.PathValue("id")
	conversationID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "conversation id không hợp lệ",
		})
		return
	}

	var req struct {
		IsPinned bool `json:"is_pinned"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
		})
		return
	}

	err = h.conversationService.PinConversation(r.Context(), conversationID, userID, req.IsPinned)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "pin conversation thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "pin conversation thành công",
	})
}
func (h *ConversationHandler) ArchiveConversation(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	idStr := r.PathValue("id")
	conversationID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "conversation id không hợp lệ",
		})
		return
	}

	var req struct {
		IsArchived bool `json:"is_archived"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
		})
		return
	}

	err = h.conversationService.ArchiveConversation(r.Context(), conversationID, userID, req.IsArchived)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "archive conversation thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "archive conversation thành công",
	})
}
func (h *ConversationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	count, err := h.conversationService.GetUnreadCount(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "lấy unread count thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Data: map[string]int64{
			"unread_count": count,
		},
	})
}
func (h *ConversationHandler) SendTypingEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	conversationIDStr := r.PathValue("id")
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "conversation id không hợp lệ",
		})
		return
	}

	var req dto.TypingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	if err := h.conversationService.SendTypingEvent(r.Context(), userID, conversationID, req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "gửi typing event thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "gửi typing event thành công",
	})
}
