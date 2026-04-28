package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserDeviceRepository struct {
	db *sql.DB
}

func NewUserDeviceRepository(db *sql.DB) *UserDeviceRepository {
	return &UserDeviceRepository{db: db}
}

func (r *UserDeviceRepository) GetByID(ctx context.Context, id int64) (*models.UserDevice, error) {
	query := `
		SELECT id, user_id, device_uuid, device_name, device_type, platform,
		       app_version, os_version, push_token, is_trusted, is_active,
		       last_seen_at, registered_at, created_at, updated_at
		FROM user_devices
		WHERE id = $1
		LIMIT 1
	`

	var d models.UserDevice
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID,
		&d.UserID,
		&d.DeviceUUID,
		&d.DeviceName,
		&d.DeviceType,
		&d.Platform,
		&d.AppVersion,
		&d.OSVersion,
		&d.PushToken,
		&d.IsTrusted,
		&d.IsActive,
		&d.LastSeenAt,
		&d.RegisteredAt,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}

func (r *UserDeviceRepository) GetByDeviceUUID(ctx context.Context, deviceUUID string) (*models.UserDevice, error) {
	query := `
		SELECT id, user_id, device_uuid, device_name, device_type, platform,
		       app_version, os_version, push_token, is_trusted, is_active,
		       last_seen_at, registered_at, created_at, updated_at
		FROM user_devices
		WHERE device_uuid = $1
		LIMIT 1
	`

	var d models.UserDevice
	err := r.db.QueryRowContext(ctx, query, deviceUUID).Scan(
		&d.ID,
		&d.UserID,
		&d.DeviceUUID,
		&d.DeviceName,
		&d.DeviceType,
		&d.Platform,
		&d.AppVersion,
		&d.OSVersion,
		&d.PushToken,
		&d.IsTrusted,
		&d.IsActive,
		&d.LastSeenAt,
		&d.RegisteredAt,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}

func (r *UserDeviceRepository) Create(ctx context.Context, d *models.UserDevice) (int64, error) {
	query := `
		INSERT INTO user_devices (
			user_id, device_uuid, device_name, device_type, platform,
			app_version, os_version, push_token, is_trusted, is_active,
			last_seen_at, registered_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14
		)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		d.UserID,
		d.DeviceUUID,
		d.DeviceName,
		d.DeviceType,
		d.Platform,
		d.AppVersion,
		d.OSVersion,
		d.PushToken,
		d.IsTrusted,
		d.IsActive,
		d.LastSeenAt,
		d.RegisteredAt,
		d.CreatedAt,
		d.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserDeviceRepository) UpdateByDeviceUUID(ctx context.Context, d *models.UserDevice) error {
	query := `
		UPDATE user_devices
		SET user_id = $1,
			device_name = $2,
			device_type = $3,
			platform = $4,
			app_version = $5,
			os_version = $6,
			push_token = $7,
			is_trusted = $8,
			is_active = $9,
			last_seen_at = $10,
			updated_at = $11
		WHERE device_uuid = $12
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		d.UserID,
		d.DeviceName,
		d.DeviceType,
		d.Platform,
		d.AppVersion,
		d.OSVersion,
		d.PushToken,
		d.IsTrusted,
		d.IsActive,
		d.LastSeenAt,
		time.Now(),
		d.DeviceUUID,
	)

	return err
}

func (r *UserDeviceRepository) UpdatePushTokenByDeviceUUID(ctx context.Context, deviceUUID string, pushToken sql.NullString) error {
	query := `
		UPDATE user_devices
		SET push_token = $1,
			updated_at = $2
		WHERE device_uuid = $3
	`

	_, err := r.db.ExecContext(ctx, query, pushToken, time.Now(), deviceUUID)
	return err
}

func (r *UserDeviceRepository) ListByUserID(ctx context.Context, userID int64) ([]models.UserDevice, error) {
	query := `
		SELECT id, user_id, device_uuid, device_name, device_type, platform,
		       app_version, os_version, push_token, is_trusted, is_active,
		       last_seen_at, registered_at, created_at, updated_at
		FROM user_devices
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.UserDevice
	for rows.Next() {
		var d models.UserDevice
		if err := rows.Scan(
			&d.ID,
			&d.UserID,
			&d.DeviceUUID,
			&d.DeviceName,
			&d.DeviceType,
			&d.Platform,
			&d.AppVersion,
			&d.OSVersion,
			&d.PushToken,
			&d.IsTrusted,
			&d.IsActive,
			&d.LastSeenAt,
			&d.RegisteredAt,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, d)
	}

	return result, rows.Err()
}
func (r *UserDeviceRepository) GetByUUIDAndUserID(ctx context.Context, uuid string, userID int64) (*models.UserDevice, error) {
	query := `
		SELECT id, user_id, device_uuid, device_name, device_type, platform,
		       app_version, os_version, push_token, is_trusted, is_active,
		       last_seen_at, registered_at, created_at, updated_at
		FROM user_devices
		WHERE device_uuid = $1 AND user_id = $2
		LIMIT 1
	`

	var d models.UserDevice
	err := r.db.QueryRowContext(ctx, query, uuid, userID).Scan(
		&d.ID,
		&d.UserID,
		&d.DeviceUUID,
		&d.DeviceName,
		&d.DeviceType,
		&d.Platform,
		&d.AppVersion,
		&d.OSVersion,
		&d.PushToken,
		&d.IsTrusted,
		&d.IsActive,
		&d.LastSeenAt,
		&d.RegisteredAt,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}
func (r *UserDeviceRepository) DisableByUUID(ctx context.Context, uuid string, userID int64) error {
	query := `
		UPDATE user_devices
		SET is_active = false, updated_at = NOW()
		WHERE device_uuid = $1 AND user_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, uuid, userID)
	return err
}
func (r *UserDeviceRepository) SetTrustedByUUID(ctx context.Context, uuid string, userID int64, trusted bool) error {
	query := `
		UPDATE user_devices
		SET is_trusted = $1, updated_at = NOW()
		WHERE device_uuid = $2 AND user_id = $3
	`
	_, err := r.db.ExecContext(ctx, query, trusted, uuid, userID)
	return err
}
