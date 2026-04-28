package middleware

import (
	"chattingapp_be/internal/dto"
	"chattingapp_be/internal/handler"
	"chattingapp_be/internal/util"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	jwtManager *util.JWTManager
}

func NewAuthMiddleware(jwtManager *util.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(dto.APIResponse{
		Success: false,
		Message: message,
	})
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeUnauthorized(w, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeUnauthorized(w, "invalid authorization header")
			return
		}

		claims, err := m.jwtManager.Verify(parts[1])
		if err != nil {
			writeUnauthorized(w, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), handler.UserIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
