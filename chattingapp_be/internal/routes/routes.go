package routes

import (
	"chattingapp_be/internal/handler"
	"chattingapp_be/internal/middleware"
	"chattingapp_be/internal/realtime"
	"net/http"
)

type Handlers struct {
	Auth         *handler.AuthHandler
	Friend       *handler.FriendHandler
	Conversation *handler.ConversationHandler
	Message      *handler.MessageHandler
	UserDevice   *handler.UserDeviceHandler
	E2EEKey      *handler.E2EEKeyHandler
	RealtimeHub  *realtime.Hub
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
	mux.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
		userID, err := realtime.ParseUserIDFromQuery(r)
		if err != nil || userID <= 0 {
			http.Error(w, "missing or invalid user_id", http.StatusBadRequest)
			return
		}

		realtime.ServeWS(handlers.RealtimeHub, w, r, userID)
	})

	// public
	mux.HandleFunc("POST /auth/register", handlers.Auth.Register)
	mux.HandleFunc("POST /auth/login", handlers.Auth.Login)
	mux.HandleFunc("POST /auth/refresh", handlers.Auth.Refresh)
	mux.HandleFunc("POST /auth/logout", handlers.Auth.Logout)
	mux.HandleFunc("POST /auth/forgot-password", handlers.Auth.ForgotPassword)
	mux.HandleFunc("POST /auth/reset-password", handlers.Auth.ResetPassword)

	// protected
	mux.Handle("GET /auth/me", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Auth.GetMyProfile)))
	mux.Handle("PATCH /auth/profile", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Auth.UpdateMyProfile)))
	mux.Handle("PATCH /auth/change-password", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Auth.ChangePassword)))

	mux.Handle("POST /friends/requests", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Friend.SendFriendRequest)))
	mux.Handle("PATCH /friends/requests/{id}", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Friend.RespondFriendRequest)))
	mux.Handle("GET /friends", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Friend.ListFriends)))
	mux.Handle("GET /friends/requests/incoming", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Friend.ListIncomingRequests)))
	mux.Handle("GET /friends/requests/outgoing", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Friend.ListOutgoingRequests)))
	mux.Handle("DELETE /friends/requests/{id}", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Friend.CancelFriendRequest)))
	mux.Handle("DELETE /friends/{user_id}", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Friend.Unfriend)))
	mux.Handle("GET /friends/search", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Friend.SearchUsers)))
	mux.Handle("GET /friends/{user_id}/mutual", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Friend.ListMutualFriends)))
	mux.Handle("POST /friends/{user_id}/block", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Friend.BlockUser)))

	mux.Handle("POST /conversations/direct", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.CreateDirectConversation)))
	mux.Handle("GET /conversations", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.ListMyConversations)))
	mux.Handle("GET /conversations/{id}", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.GetConversationDetail)))
	mux.Handle("PATCH /conversations/{id}/read", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.MarkConversationRead)))
	mux.Handle("POST /conversations/{id}/typing", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.SendTypingEvent)))
	mux.Handle("PATCH /conversations/{id}/nickname", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.UpdateMyNickname)))
	mux.Handle("PATCH /conversations/{id}/mute", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.MuteConversation)))
	mux.Handle("PATCH /conversations/{id}/pin", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.PinConversation)))
	mux.Handle("PATCH /conversations/{id}/archive", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.ArchiveConversation)))
	mux.Handle("GET /conversations/unread-count", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Conversation.GetUnreadCount)))

	mux.Handle("POST /messages", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.SendMessage)))
	mux.Handle("POST /messages/encrypted", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.SendEncryptedMessage)))
	mux.Handle("GET /messages/ciphertexts", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.ListUndeliveredCiphertexts)))
	mux.Handle("PATCH /messages/ciphertexts/{id}/delivered", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.MarkCiphertextDelivered)))
	mux.Handle("GET /messages", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.ListMessages)))
	mux.Handle("GET /messages/search", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.SearchMessages)))
	mux.Handle("GET /messages/cursor", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.ListMessagesBeforeID)))
	mux.Handle("POST /messages/{id}/forward", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.ForwardMessage)))
	mux.Handle("POST /messages/{id}/reactions", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.ReactMessage)))
	mux.Handle("GET /messages/{id}/reactions", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.ListReactions)))
	mux.Handle("DELETE /messages/{id}/reactions", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.DeleteReaction)))
	mux.Handle("PATCH /messages/{id}", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.EditMessage)))
	mux.Handle("DELETE /messages/{id}", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.DeleteMessage)))
	mux.Handle("PATCH /messages/{id}/recall", authMiddleware.RequireAuth(http.HandlerFunc(handlers.Message.RecallMessage)))

	mux.Handle("POST /devices/register", authMiddleware.RequireAuth(http.HandlerFunc(handlers.UserDevice.RegisterDevice)))
	mux.Handle("GET /devices", authMiddleware.RequireAuth(http.HandlerFunc(handlers.UserDevice.ListDevices)))
	mux.Handle("DELETE /devices/{uuid}", authMiddleware.RequireAuth(http.HandlerFunc(handlers.UserDevice.DeleteDevice)))
	mux.Handle("POST /devices/{uuid}/logout", authMiddleware.RequireAuth(http.HandlerFunc(handlers.UserDevice.LogoutDevice)))
	mux.Handle("PATCH /devices/{uuid}/disable", authMiddleware.RequireAuth(http.HandlerFunc(handlers.UserDevice.DisableDevice)))
	mux.Handle("PATCH /devices/{uuid}/trust", authMiddleware.RequireAuth(http.HandlerFunc(handlers.UserDevice.TrustDevice)))
	mux.Handle("PATCH /devices/{uuid}/untrust", authMiddleware.RequireAuth(http.HandlerFunc(handlers.UserDevice.UntrustDevice)))

	mux.Handle("POST /e2ee/keys/identity", authMiddleware.RequireAuth(http.HandlerFunc(handlers.E2EEKey.UploadIdentityKey)))
	mux.Handle("POST /e2ee/keys/signed-prekey", authMiddleware.RequireAuth(http.HandlerFunc(handlers.E2EEKey.UploadSignedPreKey)))
	mux.Handle("POST /e2ee/keys/one-time-prekeys", authMiddleware.RequireAuth(http.HandlerFunc(handlers.E2EEKey.UploadOneTimePreKeys)))
	mux.Handle("GET /e2ee/users/{user_id}/key-bundle", authMiddleware.RequireAuth(http.HandlerFunc(handlers.E2EEKey.GetUserKeyBundle)))
}
