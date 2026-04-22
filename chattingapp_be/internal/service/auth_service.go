package service

import (
	"chattingapp_be/internal/dto"
	"chattingapp_be/internal/models"
	"chattingapp_be/internal/repository"
	"chattingapp_be/internal/util"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"
	"time"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo         *repository.UserRepository
	userDeviceRepo   *repository.UserDeviceRepository
	refreshTokenRepo *repository.UserRefreshTokenRepository
	jwtManager       *util.JWTManager
	refreshDuration  time.Duration
}

func NewAuthService(
	userRepo *repository.UserRepository,
	userDeviceRepo *repository.UserDeviceRepository,
	refreshTokenRepo *repository.UserRefreshTokenRepository,
	jwtManager *util.JWTManager,
	refreshExpiresHours int,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		userDeviceRepo:   userDeviceRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
		refreshDuration:  time.Duration(refreshExpiresHours) * time.Hour,
	}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.FullName = strings.TrimSpace(req.FullName)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)

	if req.Username == "" || req.FullName == "" || req.Email == "" || req.PhoneNumber == "" || req.Password == "" {
		return nil, errors.New("thiếu thông tin đăng ký")
	}

	existingUsername, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUsername != nil {
		return nil, errors.New("username đã tồn tại")
	}

	existingEmail, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, errors.New("email đã tồn tại")
	}

	existingPhone, err := s.userRepo.GetByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if existingPhone != nil {
		return nil, errors.New("số điện thoại đã tồn tại")
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &models.User{
		Username:     req.Username,
		FullName:     req.FullName,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		PasswordHash: string(passwordHashBytes),
		Status:       "active",
		IsVerified:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	userID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = userID

	respUser := toUserResponse(user)

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}



	refreshTokenModel := &models.UserRefreshToken{
		UserID:       user.ID,
		UserDeviceID: sql.NullInt64{Valid: false},
		TokenHash:    hashToken(refreshToken),
		ExpiresAt:    now.Add(s.refreshDuration),
		RevokedAt:    sql.NullTime{Valid: false},
		LastUsedAt:   sql.NullTime{Valid: false},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
fmt.Println("LOGIN: about to create refresh token for user:", user.ID)
fmt.Println("LOGIN: token hash:", refreshTokenModel.TokenHash)
	_, err = s.refreshTokenRepo.Create(ctx, refreshTokenModel)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		User:         respUser,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	identifier := strings.TrimSpace(req.Identifier)
	if identifier == "" || req.Password == "" {
		return nil, errors.New("thiếu thông tin đăng nhập")
	}

	user, err := s.userRepo.GetByIdentifier(ctx, identifier)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("tài khoản không tồn tại")
	}

	if user.Status != "active" {
		return nil, errors.New("tài khoản không hoạt động")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("mật khẩu không đúng")
	}

	now := time.Now()

	var userDeviceID sql.NullInt64
	userDeviceID = sql.NullInt64{Valid: false}

	if strings.TrimSpace(req.DeviceUUID) != "" {
		existingDevice, err := s.userDeviceRepo.GetByDeviceUUID(ctx, strings.TrimSpace(req.DeviceUUID))
		if err != nil {
			return nil, err
		}

		device := &models.UserDevice{
			UserID:       user.ID,
			DeviceUUID:   strings.TrimSpace(req.DeviceUUID),
			DeviceName:   strings.TrimSpace(req.DeviceName),
			DeviceType:   strings.TrimSpace(req.DeviceType),
			Platform:     strings.TrimSpace(req.Platform),
			IsTrusted:    false,
			IsActive:     true,
			LastSeenAt:   sql.NullTime{Time: now, Valid: true},
			RegisteredAt: now,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if strings.TrimSpace(req.AppVersion) != "" {
			device.AppVersion = sql.NullString{String: strings.TrimSpace(req.AppVersion), Valid: true}
		}
		if strings.TrimSpace(req.OSVersion) != "" {
			device.OSVersion = sql.NullString{String: strings.TrimSpace(req.OSVersion), Valid: true}
		}

		if existingDevice == nil {
			deviceID, err := s.userDeviceRepo.Create(ctx, device)
			if err != nil {
				return nil, err
			}
			userDeviceID = sql.NullInt64{Int64: deviceID, Valid: true}
		} else {
			err := s.userDeviceRepo.UpdateByDeviceUUID(ctx, device)
			if err != nil {
				return nil, err
			}
			userDeviceID = sql.NullInt64{Int64: existingDevice.ID, Valid: true}
		}
	}

	respUser := toUserResponse(user)

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshTokenModel := &models.UserRefreshToken{
		UserID:       user.ID,
		UserDeviceID: userDeviceID,
		TokenHash:    hashToken(refreshToken),
		ExpiresAt:    now.Add(s.refreshDuration),
		RevokedAt:    sql.NullTime{Valid: false},
		LastUsedAt:   sql.NullTime{Valid: false},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err = s.refreshTokenRepo.Create(ctx, refreshTokenModel)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		User:         respUser,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.AuthResponse, error) {
	req.RefreshToken = strings.TrimSpace(req.RefreshToken)
	if req.RefreshToken == "" {
		return nil, errors.New("thiếu refresh token")
	}

	claims, err := s.jwtManager.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("refresh token không hợp lệ hoặc đã hết hạn")
	}

	storedToken, err := s.refreshTokenRepo.GetByTokenHash(ctx, hashToken(req.RefreshToken))
	if err != nil {
		return nil, err
	}
	if storedToken == nil {
		return nil, errors.New("refresh token không tồn tại")
	}

	if storedToken.RevokedAt.Valid {
		return nil, errors.New("refresh token đã bị thu hồi")
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return nil, errors.New("refresh token đã hết hạn")
	}

	if storedToken.UserID != claims.UserID {
		return nil, errors.New("refresh token không khớp user")
	}

	now := time.Now()

	err = s.refreshTokenRepo.UpdateLastUsedAt(ctx, storedToken.ID, sql.NullTime{
		Time:  now,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}

	err = s.refreshTokenRepo.RevokeByID(ctx, storedToken.ID, sql.NullTime{
		Time:  now,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}

	newAccessToken, err := s.jwtManager.GenerateAccessToken(storedToken.UserID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(storedToken.UserID)
	if err != nil {
		return nil, err
	}

	newStoredToken := &models.UserRefreshToken{
		UserID:       storedToken.UserID,
		UserDeviceID: storedToken.UserDeviceID,
		TokenHash:    hashToken(newRefreshToken),
		ExpiresAt:    now.Add(s.refreshDuration),
		RevokedAt:    sql.NullTime{Valid: false},
		LastUsedAt:   sql.NullTime{Valid: false},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err = s.refreshTokenRepo.Create(ctx, newStoredToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("không tìm thấy user")
	}

	respUser := toUserResponse(user)

	return &dto.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		User:         respUser,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req dto.LogoutRequest) error {
	req.RefreshToken = strings.TrimSpace(req.RefreshToken)
	if req.RefreshToken == "" {
		return errors.New("thiếu refresh token")
	}

	storedToken, err := s.refreshTokenRepo.GetByTokenHash(ctx, hashToken(req.RefreshToken))
	if err != nil {
		return err
	}
	if storedToken == nil {
		return nil
	}

	if storedToken.RevokedAt.Valid {
		return nil
	}

	return s.refreshTokenRepo.RevokeByID(ctx, storedToken.ID, sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	})
}

func (s *AuthService) GetMyProfile(ctx context.Context, userID int64) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("không tìm thấy user")
	}

	resp := toUserResponse(user)
	return &resp, nil
}

func toUserResponse(u *models.User) dto.UserResponse {
	var avatarURL *string
	if u.AvatarURL.Valid {
		avatarURL = &u.AvatarURL.String
	}

	var bio *string
	if u.Bio.Valid {
		bio = &u.Bio.String
	}

	var lastSeenAt *string
	if u.LastSeenAt.Valid {
		v := u.LastSeenAt.Time.Format(time.RFC3339)
		lastSeenAt = &v
	}

	return dto.UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		FullName:    u.FullName,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		AvatarURL:   avatarURL,
		Bio:         bio,
		Status:      u.Status,
		IsVerified:  u.IsVerified,
		LastSeenAt:  lastSeenAt,
		CreatedAt:   u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   u.UpdatedAt.Format(time.RFC3339),
	}
}

func nullStringPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}