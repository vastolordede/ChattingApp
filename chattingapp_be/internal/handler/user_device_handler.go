package handler

import (
	"chattingapp_be/internal/dto"
	"chattingapp_be/internal/service"
	"encoding/json"
	"net/http"
)

type UserDeviceHandler struct {
	deviceService *service.UserDeviceService
}

func NewUserDeviceHandler(deviceService *service.UserDeviceService) *UserDeviceHandler {
	return &UserDeviceHandler{
		deviceService: deviceService,
	}
}

func (h *UserDeviceHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	var req dto.RegisterDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "request body không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	if err := h.deviceService.RegisterOrUpdateDevice(r.Context(), userID, req); err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "đăng ký device thất bại",
			Error:   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "đăng ký device thành công",
	})
}
func (h *UserDeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "unauthorized"})
		return
	}

	devices, err := h.deviceService.ListMyDevices(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    devices,
	})
}
func (h *UserDeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "unauthorized"})
		return
	}

	uuid := r.PathValue("uuid")

	if err := h.deviceService.DeleteDevice(r.Context(), userID, uuid); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "device deleted",
	})
}
func (h *UserDeviceHandler) LogoutDevice(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "unauthorized"})
		return
	}

	uuid := r.PathValue("uuid")

	if err := h.deviceService.LogoutDevice(r.Context(), userID, uuid); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "logout device success",
	})
}
func (h *UserDeviceHandler) DisableDevice(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "unauthorized"})
		return
	}

	uuid := r.PathValue("uuid")

	if err := h.deviceService.DisableDevice(r.Context(), userID, uuid); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "device disabled",
	})
}

func (h *UserDeviceHandler) TrustDevice(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "unauthorized"})
		return
	}

	uuid := r.PathValue("uuid")

	if err := h.deviceService.TrustDevice(r.Context(), userID, uuid); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "device trusted",
	})
}

func (h *UserDeviceHandler) UntrustDevice(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "unauthorized"})
		return
	}

	uuid := r.PathValue("uuid")

	if err := h.deviceService.UntrustDevice(r.Context(), userID, uuid); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "device untrusted",
	})
}
