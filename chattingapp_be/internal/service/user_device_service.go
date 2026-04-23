package service

import (
	"chattingapp_be/internal/dto"
	"chattingapp_be/internal/models"
	"chattingapp_be/internal/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserDeviceService struct {
	repo *repository.UserDeviceRepository
	tokenRepo  *repository.UserRefreshTokenRepository
}

func NewUserDeviceService(
	repo *repository.UserDeviceRepository,
	tokenRepo *repository.UserRefreshTokenRepository,
) *UserDeviceService {
	return &UserDeviceService{
		repo:      repo,
		tokenRepo: tokenRepo,
	}
}

func (s *UserDeviceService) RegisterOrUpdateDevice(
	ctx context.Context,
	userID int64,
	req dto.RegisterDeviceRequest,
) error {
	now := time.Now()

	existing, err := s.repo.GetByDeviceUUID(ctx, req.DeviceUUID)
	if err != nil {
		return err
	}

	device := &models.UserDevice{
		UserID:       userID,
		DeviceUUID:   req.DeviceUUID,
		DeviceName:   req.DeviceName,
		DeviceType:   req.DeviceType,
		Platform:     req.Platform,
		IsTrusted:    false,
		IsActive:     true,
		LastSeenAt:   sql.NullTime{Time: now, Valid: true},
		RegisteredAt: now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if req.AppVersion != "" {
		device.AppVersion = sql.NullString{String: req.AppVersion, Valid: true}
	}
	if req.OSVersion != "" {
		device.OSVersion = sql.NullString{String: req.OSVersion, Valid: true}
	}
	if req.PushToken != "" {
		device.PushToken = sql.NullString{String: req.PushToken, Valid: true}
	}

	if existing == nil {
		_, err := s.repo.Create(ctx, device)
		return err
	}

	return s.repo.UpdateByDeviceUUID(ctx, device)
}

func (s *UserDeviceService) ListMyDevices(ctx context.Context, userID int64) ([]models.UserDevice, error) {
	return s.repo.ListByUserID(ctx, userID)
}
func (s *UserDeviceService) DeleteDevice(ctx context.Context, userID int64, uuid string) error {
	device, err := s.repo.GetByUUIDAndUserID(ctx, uuid, userID)
	if err != nil {
		return err
	}
	if device == nil {
		return errors.New("device not found")
	}

	if err := s.repo.DisableByUUID(ctx, uuid, userID); err != nil {
		return err
	}

	return s.tokenRepo.RevokeByDeviceID(ctx, device.ID)
}

func (s *UserDeviceService) LogoutDevice(ctx context.Context, userID int64, uuid string) error {
	device, err := s.repo.GetByUUIDAndUserID(ctx, uuid, userID)
	if err != nil {
		return err
	}
	if device == nil {
		return errors.New("device not found")
	}

	return s.tokenRepo.RevokeByDeviceID(ctx, device.ID)
}
func (s *UserDeviceService) DisableDevice(ctx context.Context, userID int64, uuid string) error {
	device, err := s.repo.GetByUUIDAndUserID(ctx, uuid, userID)
	if err != nil {
		return err
	}
	if device == nil {
		return errors.New("device not found")
	}

	if err := s.repo.DisableByUUID(ctx, uuid, userID); err != nil {
		return err
	}

	return s.tokenRepo.RevokeByDeviceID(ctx, device.ID)
}

func (s *UserDeviceService) TrustDevice(ctx context.Context, userID int64, uuid string) error {
	device, err := s.repo.GetByUUIDAndUserID(ctx, uuid, userID)
	if err != nil {
		return err
	}
	if device == nil {
		return errors.New("device not found")
	}

	return s.repo.SetTrustedByUUID(ctx, uuid, userID, true)
}

func (s *UserDeviceService) UntrustDevice(ctx context.Context, userID int64, uuid string) error {
	device, err := s.repo.GetByUUIDAndUserID(ctx, uuid, userID)
	if err != nil {
		return err
	}
	if device == nil {
		return errors.New("device not found")
	}

	return s.repo.SetTrustedByUUID(ctx, uuid, userID, false)
}