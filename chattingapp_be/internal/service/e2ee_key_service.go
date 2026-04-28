package service

import (
	"chattingapp_be/internal/dto"
	"chattingapp_be/internal/models"
	"chattingapp_be/internal/repository"
	"context"
	"database/sql"
	"errors"
	"strings"
)

type E2EEKeyService struct {
	userDeviceRepo    *repository.UserDeviceRepository
	identityKeyRepo   *repository.DeviceIdentityKeyRepository
	signedPreKeyRepo  *repository.DeviceSignedPreKeyRepository
	oneTimePreKeyRepo *repository.DeviceOneTimePreKeyRepository
}

func NewE2EEKeyService(
	userDeviceRepo *repository.UserDeviceRepository,
	identityKeyRepo *repository.DeviceIdentityKeyRepository,
	signedPreKeyRepo *repository.DeviceSignedPreKeyRepository,
	oneTimePreKeyRepo *repository.DeviceOneTimePreKeyRepository,
) *E2EEKeyService {
	return &E2EEKeyService{
		userDeviceRepo:    userDeviceRepo,
		identityKeyRepo:   identityKeyRepo,
		signedPreKeyRepo:  signedPreKeyRepo,
		oneTimePreKeyRepo: oneTimePreKeyRepo,
	}
}

func (s *E2EEKeyService) UploadIdentityKey(
	ctx context.Context,
	userID int64,
	req dto.UploadIdentityKeyRequest,
) (int64, error) {
	device, err := s.getOwnedActiveDevice(ctx, userID, req.DeviceUUID)
	if err != nil {
		return 0, err
	}

	if strings.TrimSpace(req.PublicKey) == "" {
		return 0, errors.New("public_key is required")
	}
	if strings.TrimSpace(req.Fingerprint) == "" {
		return 0, errors.New("fingerprint is required")
	}

	algorithm := strings.TrimSpace(req.Algorithm)
	if algorithm == "" {
		algorithm = "X25519"
	}

	version := req.Version
	if version <= 0 {
		version = 1
	}

	key := &models.DeviceIdentityKey{
		DeviceID:    device.ID,
		PublicKey:   strings.TrimSpace(req.PublicKey),
		Algorithm:   algorithm,
		Fingerprint: strings.TrimSpace(req.Fingerprint),
		Version:     version,
		IsActive:    true,
		ExpiredAt:   sql.NullTime{},
		RevokedAt:   sql.NullTime{},
	}

	return s.identityKeyRepo.ReplaceActive(ctx, key)
}

func (s *E2EEKeyService) UploadSignedPreKey(
	ctx context.Context,
	userID int64,
	req dto.UploadSignedPreKeyRequest,
) (int64, error) {
	device, err := s.getOwnedActiveDevice(ctx, userID, req.DeviceUUID)
	if err != nil {
		return 0, err
	}

	if req.KeyID <= 0 {
		return 0, errors.New("key_id must be greater than 0")
	}
	if strings.TrimSpace(req.PublicKey) == "" {
		return 0, errors.New("public_key is required")
	}
	if strings.TrimSpace(req.Signature) == "" {
		return 0, errors.New("signature is required")
	}

	algorithm := strings.TrimSpace(req.Algorithm)
	if algorithm == "" {
		algorithm = "X25519"
	}

	version := req.Version
	if version <= 0 {
		version = 1
	}

	key := &models.DeviceSignedPreKey{
		DeviceID:  device.ID,
		KeyID:     req.KeyID,
		PublicKey: strings.TrimSpace(req.PublicKey),
		Signature: strings.TrimSpace(req.Signature),
		Algorithm: algorithm,
		Version:   version,
		IsActive:  true,
		ExpiredAt: sql.NullTime{},
		RevokedAt: sql.NullTime{},
	}

	return s.signedPreKeyRepo.ReplaceActive(ctx, key)
}

func (s *E2EEKeyService) UploadOneTimePreKeys(
	ctx context.Context,
	userID int64,
	req dto.UploadOneTimePreKeysRequest,
) (int, error) {
	device, err := s.getOwnedActiveDevice(ctx, userID, req.DeviceUUID)
	if err != nil {
		return 0, err
	}

	if len(req.PreKeys) == 0 {
		return 0, errors.New("prekeys is required")
	}

	if len(req.PreKeys) > 100 {
		return 0, errors.New("maximum 100 one-time prekeys per upload")
	}

	keys := make([]models.DeviceOneTimePreKey, 0, len(req.PreKeys))
	seen := make(map[int]bool)

	for _, item := range req.PreKeys {
		if item.KeyID <= 0 {
			return 0, errors.New("prekey key_id must be greater than 0")
		}
		if seen[item.KeyID] {
			return 0, errors.New("duplicated prekey key_id in request")
		}
		seen[item.KeyID] = true

		if strings.TrimSpace(item.PublicKey) == "" {
			return 0, errors.New("prekey public_key is required")
		}

		algorithm := strings.TrimSpace(item.Algorithm)
		if algorithm == "" {
			algorithm = "X25519"
		}

		version := item.Version
		if version <= 0 {
			version = 1
		}

		keys = append(keys, models.DeviceOneTimePreKey{
			DeviceID:  device.ID,
			KeyID:     item.KeyID,
			PublicKey: strings.TrimSpace(item.PublicKey),
			Algorithm: algorithm,
			Version:   version,
			IsUsed:    false,
			UsedAt:    sql.NullTime{},
			ExpiredAt: sql.NullTime{},
		})
	}

	return s.oneTimePreKeyRepo.CreateBatch(ctx, keys)
}

func (s *E2EEKeyService) GetUserKeyBundle(
	ctx context.Context,
	targetUserID int64,
) (*dto.UserKeyBundleResponse, error) {
	if targetUserID <= 0 {
		return nil, errors.New("target user id is invalid")
	}

	devices, err := s.userDeviceRepo.ListByUserID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	resp := &dto.UserKeyBundleResponse{
		UserID:  targetUserID,
		Devices: []dto.DeviceKeyBundleResponse{},
	}

	for _, device := range devices {
		if !device.IsActive {
			continue
		}

		identityKey, err := s.identityKeyRepo.GetActiveByDeviceID(ctx, device.ID)
		if err != nil {
			return nil, err
		}
		if identityKey == nil {
			continue
		}

		signedPreKey, err := s.signedPreKeyRepo.GetActiveByDeviceID(ctx, device.ID)
		if err != nil {
			return nil, err
		}
		if signedPreKey == nil {
			continue
		}

		oneTimePreKey, err := s.oneTimePreKeyRepo.ConsumeOneByDeviceID(ctx, device.ID)
		if err != nil {
			return nil, err
		}

		deviceBundle := dto.DeviceKeyBundleResponse{
			DeviceID:   device.ID,
			DeviceUUID: device.DeviceUUID,
			IdentityKey: dto.IdentityKeyResponse{
				KeyID:       identityKey.ID,
				DeviceID:    identityKey.DeviceID,
				PublicKey:   identityKey.PublicKey,
				Algorithm:   identityKey.Algorithm,
				Fingerprint: identityKey.Fingerprint,
				Version:     identityKey.Version,
			},
			SignedPreKey: dto.SignedPreKeyResponse{
				KeyID:     signedPreKey.ID,
				DeviceID:  signedPreKey.DeviceID,
				PreKeyID:  signedPreKey.KeyID,
				PublicKey: signedPreKey.PublicKey,
				Signature: signedPreKey.Signature,
				Algorithm: signedPreKey.Algorithm,
				Version:   signedPreKey.Version,
			},
			HasOneTimePreKey: false,
		}

		if oneTimePreKey != nil {
			deviceBundle.OneTimePreKey = &dto.OneTimePreKeyResponse{
				KeyID:     oneTimePreKey.ID,
				DeviceID:  oneTimePreKey.DeviceID,
				PreKeyID:  oneTimePreKey.KeyID,
				PublicKey: oneTimePreKey.PublicKey,
				Algorithm: oneTimePreKey.Algorithm,
				Version:   oneTimePreKey.Version,
			}
			deviceBundle.HasOneTimePreKey = true
		}

		resp.Devices = append(resp.Devices, deviceBundle)
	}

	return resp, nil
}

func (s *E2EEKeyService) getOwnedActiveDevice(ctx context.Context, userID int64, deviceUUID string) (*models.UserDevice, error) {
	deviceUUID = strings.TrimSpace(deviceUUID)
	if deviceUUID == "" {
		return nil, errors.New("device_uuid is required")
	}

	device, err := s.userDeviceRepo.GetByUUIDAndUserID(ctx, deviceUUID, userID)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, errors.New("device not found or not owned by current user")
	}
	if !device.IsActive {
		return nil, errors.New("device is inactive")
	}

	return device, nil
}
