// @title ChattingApp Backend API
// @version 1.0
// @description API tài liệu cho ChattingApp Backend.
// @description Protected APIs dùng Authorization: Bearer <access_token>
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Nhập theo dạng: Bearer <access_token>
package main

import (
	_ "chattingapp_be/docs"
	"chattingapp_be/internal/config"
	"chattingapp_be/internal/database"
	"chattingapp_be/internal/handler"
	"chattingapp_be/internal/middleware"
	"chattingapp_be/internal/realtime"
	"chattingapp_be/internal/repository"
	"chattingapp_be/internal/routes"
	"chattingapp_be/internal/service"
	_ "chattingapp_be/internal/swaggerdocs"
	"chattingapp_be/internal/util"
	"fmt"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config thất bại: %v", err)
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("kết nối database thất bại: %v", err)
	}
	defer db.Close()

	jwtManager := util.NewJWTManager(
		cfg.JWTSecret,
		cfg.JWTExpiresHours,
		cfg.RefreshExpiresHours,
	)

	userRepo := repository.NewUserRepository(db)
	friendRequestRepo := repository.NewFriendRequestRepository(db)
	friendshipRepo := repository.NewFriendshipRepository(db)
	conversationRepo := repository.NewConversationRepository(db)
	conversationMemberRepo := repository.NewConversationMemberRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	messageAttachmentRepo := repository.NewMessageAttachmentRepository(db)
	userDeviceRepo := repository.NewUserDeviceRepository(db)
	userRefreshTokenRepo := repository.NewUserRefreshTokenRepository(db)
	passwordResetTokenRepo := repository.NewPasswordResetTokenRepository(db)
	realtimeHub := realtime.NewHub()

	deviceIdentityKeyRepo := repository.NewDeviceIdentityKeyRepository(db)
	deviceSignedPreKeyRepo := repository.NewDeviceSignedPreKeyRepository(db)
	deviceOneTimePreKeyRepo := repository.NewDeviceOneTimePreKeyRepository(db)
	messageCiphertextRepo := repository.NewMessageCiphertextRepository(db)

	authService := service.NewAuthService(
		userRepo,
		userDeviceRepo,
		userRefreshTokenRepo,
		passwordResetTokenRepo,
		jwtManager,
		cfg.RefreshExpiresHours,
	)

	friendService := service.NewFriendService(
		db,
		friendRequestRepo,
		friendshipRepo,
		userRepo,
	)

	conversationService := service.NewConversationService(
		db,
		conversationRepo,
		conversationMemberRepo,
		userRepo,
		messageRepo,
		realtimeHub,
	)

	messageService := service.NewMessageService(
		db,
		messageRepo,
		messageAttachmentRepo,
		conversationRepo,
		conversationMemberRepo,
		userRepo,
		messageCiphertextRepo,
		userDeviceRepo,
		realtimeHub,
	)

	userDeviceService := service.NewUserDeviceService(userDeviceRepo, userRefreshTokenRepo)

	e2eeKeyService := service.NewE2EEKeyService(
		userDeviceRepo,
		deviceIdentityKeyRepo,
		deviceSignedPreKeyRepo,
		deviceOneTimePreKeyRepo,
	)

	authHandler := handler.NewAuthHandler(authService)
	friendHandler := handler.NewFriendHandler(friendService)
	conversationHandler := handler.NewConversationHandler(conversationService)
	messageHandler := handler.NewMessageHandler(messageService)
	userDeviceHandler := handler.NewUserDeviceHandler(userDeviceService)

	e2eeKeyHandler := handler.NewE2EEKeyHandler(e2eeKeyService)

	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	mux := http.NewServeMux()
	routes.RegisterRoutes(
		mux,
		routes.Handlers{
			Auth:         authHandler,
			Friend:       friendHandler,
			Conversation: conversationHandler,
			Message:      messageHandler,
			UserDevice:   userDeviceHandler,
			E2EEKey:      e2eeKeyHandler,
			RealtimeHub:  realtimeHub,
		},
		authMiddleware,
	)

	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	addr := ":" + cfg.AppPort
	fmt.Printf("server đang chạy tại http://localhost%s\n", addr)
	fmt.Printf("swagger đang chạy tại http://localhost%s/swagger/index.html\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("khởi động server thất bại: %v", err)
	}
}
