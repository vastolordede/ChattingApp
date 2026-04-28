package handler

import (
	"chattingapp_be/internal/dto"
	"chattingapp_be/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type E2EEKeyHandler struct {
	keyService *service.E2EEKeyService
}

func NewE2EEKeyHandler(keyService *service.E2EEKeyService) *E2EEKeyHandler {
	return &E2EEKeyHandler{keyService: keyService}
}

func (h *E2EEKeyHandler) UploadIdentityKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
			Error:   "user id not found in context",
		})
		return
	}

	var req dto.UploadIdentityKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request body",
			Error:   err.Error(),
		})
		return
	}

	id, err := h.keyService.UploadIdentityKey(r.Context(), userID, req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "upload identity key thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "upload identity key thành công",
		Data: dto.E2EEKeyUploadResponse{
			Message: "identity key uploaded successfully",
			ID:      id,
		},
	})
}

func (h *E2EEKeyHandler) UploadSignedPreKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
			Error:   "user id not found in context",
		})
		return
	}

	var req dto.UploadSignedPreKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request body",
			Error:   err.Error(),
		})
		return
	}

	id, err := h.keyService.UploadSignedPreKey(r.Context(), userID, req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "upload signed prekey thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "upload signed prekey thành công",
		Data: dto.E2EEKeyUploadResponse{
			Message: "signed prekey uploaded successfully",
			ID:      id,
		},
	})
}

func (h *E2EEKeyHandler) UploadOneTimePreKeys(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
			Error:   "user id not found in context",
		})
		return
	}

	var req dto.UploadOneTimePreKeysRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request body",
			Error:   err.Error(),
		})
		return
	}

	count, err := h.keyService.UploadOneTimePreKeys(r.Context(), userID, req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "upload one-time prekeys thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "upload one-time prekeys thành công",
		Data: dto.E2EEKeyUploadResponse{
			Message: "one-time prekeys uploaded successfully",
			Count:   count,
		},
	})
}

func (h *E2EEKeyHandler) GetUserKeyBundle(w http.ResponseWriter, r *http.Request) {
	targetUserID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || targetUserID <= 0 {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid user id",
			Error:   "user_id must be a positive integer",
		})
		return
	}

	resp, err := h.keyService.GetUserKeyBundle(r.Context(), targetUserID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "lấy key bundle thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "lấy key bundle thành công",
		Data:    resp,
	})
}
