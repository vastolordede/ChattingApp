package handler

import (
	"chattingapp_be/internal/dto"
	"context"
	"encoding/json"
	"net/http"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

func writeJSON(w http.ResponseWriter, status int, resp dto.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

func getUserIDFromContext(r *http.Request) (int64, bool) {
	v := r.Context().Value(UserIDContextKey)
	if v == nil {
		return 0, false
	}
	userID, ok := v.(int64)
	return userID, ok
}

func withUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}
