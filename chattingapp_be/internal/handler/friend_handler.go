package handler

import (
	"chattingapp_be/internal/dto"
	"chattingapp_be/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type FriendHandler struct {
	friendService *service.FriendService
}

func NewFriendHandler(friendService *service.FriendService) *FriendHandler {
	return &FriendHandler{
		friendService: friendService,
	}
}

func (h *FriendHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	var req dto.SendFriendRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	if err := h.friendService.SendFriendRequest(r.Context(), userID, req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "gửi lời mời kết bạn thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "gửi lời mời kết bạn thành công",
	})
}

func (h *FriendHandler) RespondFriendRequest(w http.ResponseWriter, r *http.Request) {
	friendRequestIDStr := r.PathValue("id")
	friendRequestID, err := strconv.ParseInt(friendRequestIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "friend request id không hợp lệ",
		})
		return
	}

	var req dto.RespondFriendRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	if err := h.friendService.RespondFriendRequest(r.Context(), userID, friendRequestID, req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "xử lý lời mời kết bạn thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "xử lý lời mời kết bạn thành công",
	})
}

func (h *FriendHandler) ListFriends(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	items, err := h.friendService.ListFriends(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "lấy danh sách bạn bè thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "lấy danh sách bạn bè thành công",
		Data:    items,
	})
}
func (h *FriendHandler) ListIncomingRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	items, err := h.friendService.ListIncomingRequests(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "lấy danh sách lời mời đến thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "lấy danh sách lời mời đến thành công",
		Data:    items,
	})
}

func (h *FriendHandler) ListOutgoingRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	items, err := h.friendService.ListOutgoingRequests(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "lấy danh sách lời mời đi thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "lấy danh sách lời mời đi thành công",
		Data:    items,
	})
}

func (h *FriendHandler) CancelFriendRequest(w http.ResponseWriter, r *http.Request) {
	friendRequestIDStr := r.PathValue("id")
	friendRequestID, err := strconv.ParseInt(friendRequestIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "friend request id không hợp lệ",
		})
		return
	}

	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	if err := h.friendService.CancelFriendRequest(r.Context(), userID, friendRequestID); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "hủy lời mời kết bạn thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "hủy lời mời kết bạn thành công",
	})
}

func (h *FriendHandler) Unfriend(w http.ResponseWriter, r *http.Request) {
	targetUserIDStr := r.PathValue("user_id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "user id không hợp lệ",
		})
		return
	}

	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	if err := h.friendService.Unfriend(r.Context(), userID, targetUserID); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "unfriend thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "unfriend thành công",
	})
}
func (h *FriendHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	q := r.URL.Query().Get("q")

	items, err := h.friendService.SearchUsers(r.Context(), userID, q)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "tìm user thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "tìm user thành công",
		Data:    items,
	})
}

func (h *FriendHandler) ListMutualFriends(w http.ResponseWriter, r *http.Request) {
	targetUserIDStr := r.PathValue("user_id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "user id không hợp lệ",
		})
		return
	}

	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	items, err := h.friendService.ListMutualFriends(r.Context(), userID, targetUserID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "lấy mutual friends thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "lấy mutual friends thành công",
		Data:    items,
	})
}

func (h *FriendHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	targetUserIDStr := r.PathValue("user_id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "user id không hợp lệ",
		})
		return
	}

	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	if err := h.friendService.BlockUser(r.Context(), userID, targetUserID); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "block user thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "block user thành công",
	})
}
