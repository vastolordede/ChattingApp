package handler

import (
	"chattingapp_be/internal/dto"
	"chattingapp_be/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	var req dto.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.messageService.SendMessage(r.Context(), userID, req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "gửi tin nhắn thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "gửi tin nhắn thành công",
		Data:    resp,
	})
}

func (h *MessageHandler) EditMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	messageIDStr := r.PathValue("id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "message id không hợp lệ",
		})
		return
	}

	var req dto.EditMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.messageService.EditMessage(r.Context(), userID, messageID, req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "sửa tin nhắn thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "sửa tin nhắn thành công",
		Data:    resp,
	})
}

func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	messageIDStr := r.PathValue("id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "message id không hợp lệ",
		})
		return
	}

	var req dto.DeleteMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.SoftDelete = true
	}

	if err := h.messageService.DeleteMessage(r.Context(), userID, messageID, req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "xóa tin nhắn thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "xóa tin nhắn thành công",
	})
}

func (h *MessageHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	conversationIDStr := r.URL.Query().Get("conversation_id")
	if conversationIDStr == "" {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "thiếu conversation_id",
		})
		return
	}

	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "conversation id không hợp lệ",
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

	items, err := h.messageService.ListMessages(r.Context(), userID, conversationID, page, limit)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "lấy danh sách tin nhắn thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "lấy danh sách tin nhắn thành công",
		Data:    items,
	})
}
func (h *MessageHandler) SearchMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	conversationIDStr := r.URL.Query().Get("conversation_id")
	if conversationIDStr == "" {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "thiếu conversation_id",
		})
		return
	}

	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "conversation id không hợp lệ",
		})
		return
	}

	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "thiếu từ khóa tìm kiếm q",
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

	items, err := h.messageService.SearchMessages(r.Context(), userID, conversationID, keyword, page, limit)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "tìm kiếm tin nhắn thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "tìm kiếm tin nhắn thành công",
		Data:    items,
	})
}

func (h *MessageHandler) ListMessagesBeforeID(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	conversationIDStr := r.URL.Query().Get("conversation_id")
	if conversationIDStr == "" {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "thiếu conversation_id",
		})
		return
	}

	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "conversation id không hợp lệ",
		})
		return
	}

	beforeIDStr := r.URL.Query().Get("before_id")
	if beforeIDStr == "" {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "thiếu before_id",
		})
		return
	}

	beforeID, err := strconv.ParseInt(beforeIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "before_id không hợp lệ",
		})
		return
	}

	limit := 20
	if v := r.URL.Query().Get("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 {
			limit = l
		}
	}

	items, err := h.messageService.ListMessagesBeforeID(r.Context(), userID, conversationID, beforeID, limit)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "lấy tin nhắn theo cursor thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "lấy tin nhắn theo cursor thành công",
		Data:    items,
	})
}

func (h *MessageHandler) ForwardMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	messageIDStr := r.PathValue("id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "message id không hợp lệ",
		})
		return
	}

	var req dto.ForwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.messageService.ForwardMessage(r.Context(), userID, messageID, req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "forward tin nhắn thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "forward tin nhắn thành công",
		Data:    resp,
	})
}

func (h *MessageHandler) ReactMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	messageIDStr := r.PathValue("id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "message id không hợp lệ",
		})
		return
	}

	var req dto.ReactMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.messageService.ReactMessage(r.Context(), userID, messageID, req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "reaction tin nhắn thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "reaction tin nhắn thành công",
		Data:    resp,
	})
}

func (h *MessageHandler) DeleteReaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	messageIDStr := r.PathValue("id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "message id không hợp lệ",
		})
		return
	}

	if err := h.messageService.DeleteReaction(r.Context(), userID, messageID); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "xóa reaction thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "xóa reaction thành công",
	})
}

func (h *MessageHandler) ListReactions(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	messageIDStr := r.PathValue("id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "message id không hợp lệ",
		})
		return
	}

	items, err := h.messageService.ListReactions(r.Context(), userID, messageID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "lấy reaction thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "lấy reaction thành công",
		Data:    items,
	})
}
func (h *MessageHandler) RecallMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	messageIDStr := r.PathValue("id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "message id không hợp lệ",
		})
		return
	}

	resp, err := h.messageService.RecallMessage(r.Context(), userID, messageID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "thu hồi tin nhắn thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "thu hồi tin nhắn thành công",
		Data:    resp,
	})
}
func (h *MessageHandler) SendEncryptedMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	var req dto.SendEncryptedMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.messageService.SendEncryptedMessage(r.Context(), userID, req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "gửi encrypted message thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "gửi encrypted message thành công",
		Data:    resp,
	})
}

func (h *MessageHandler) ListUndeliveredCiphertexts(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	deviceUUID := r.URL.Query().Get("device_uuid")
	items, err := h.messageService.ListUndeliveredCiphertextsForDevice(r.Context(), userID, deviceUUID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "lấy ciphertext thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "lấy ciphertext thành công",
		Data:    items,
	})
}

func (h *MessageHandler) MarkCiphertextDelivered(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	idStr := r.PathValue("id")
	ciphertextID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || ciphertextID <= 0 {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "ciphertext id không hợp lệ",
		})
		return
	}

	if err := h.messageService.MarkCiphertextDelivered(r.Context(), userID, ciphertextID); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "mark ciphertext delivered thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "mark ciphertext delivered thành công",
	})
}
