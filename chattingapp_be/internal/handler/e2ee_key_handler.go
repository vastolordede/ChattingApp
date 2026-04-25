package handler

import (
	"chattingapp_be/internal/dto"
	"chattingapp_be/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
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
	targetUserID, err := parseUserIDFromKeyBundlePath(r.URL.Path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request path",
			Error:   err.Error(),
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

func parseUserIDFromKeyBundlePath(path string) (int64, error) {
	// Expected path: /e2ee/users/{user_id}/key-bundle
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 4 {
		return 0, errors.New("invalid key bundle path")
	}

	if parts[0] != "e2ee" || parts[1] != "users" || parts[3] != "key-bundle" {
		return 0, errors.New("invalid key bundle path")
	}

	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid user id")
	}

	return id, nil
}