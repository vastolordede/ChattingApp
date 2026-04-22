package routes

import (
	"chattingapp_be/internal/handler"
	"chattingapp_be/internal/middleware"
	"net/http"
)

type Handlers struct {
	Auth         *handler.AuthHandler
	Friend       *handler.FriendHandler
	Conversation *handler.ConversationHandler
	Message      *handler.MessageHandler
	UserDevice   *handler.UserDeviceHandler
}

func RegisterRoutes(
	mux *http.ServeMux,
	handlers Handlers,
	authMiddleware *middleware.AuthMiddleware,
) {
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"pong"}`))
	})

	// public
	mux.HandleFunc("POST /auth/register", handlers.Auth.Register)
	mux.HandleFunc("POST /auth/login", handlers.Auth.Login)
	mux.HandleFunc("POST /auth/refresh", handlers.Auth.Refresh)
	mux.HandleFunc("POST /auth/logout", handlers.Auth.Logout)

	// protected
	mux.Handle("GET /auth/me", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Auth.GetMyProfile)))
	mux.Handle("PATCH /auth/profile", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Auth.UpdateMyProfile)))
	mux.Handle("PATCH /auth/change-password", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Auth.ChangePassword)))

	mux.Handle("POST /friends/requests", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Friend.SendFriendRequest)))
	mux.Handle("PATCH /friends/requests/{id}", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Friend.RespondFriendRequest)))
	mux.Handle("GET /friends", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Friend.ListFriends)))

	mux.Handle("POST /conversations/direct", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.CreateDirectConversation)))
	mux.Handle("GET /conversations", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.ListMyConversations)))
	mux.Handle("GET /conversations/{id}", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.GetConversationDetail)))
	mux.Handle("PATCH /conversations/{id}/read", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.MarkConversationRead)))
	mux.Handle("PATCH /conversations/{id}/nickname", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.UpdateMyNickname)))
	mux.Handle("PATCH /conversations/{id}/mute", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.MuteConversation)))

	mux.Handle("POST /messages", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.SendMessage)))
	mux.Handle("PATCH /messages/{id}", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.EditMessage)))
	mux.Handle("DELETE /messages/{id}", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.DeleteMessage)))
	mux.Handle("GET /messages", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.ListMessages)))

	mux.Handle("POST /devices/register", authMiddleware.RequireAuth(http.HandlerFunc(handlers.UserDevice.RegisterDevice)))
}